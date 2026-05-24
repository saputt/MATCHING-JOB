package scraper

import (
	"context"
	"fmt"
	"job-matching-scraper/internal/entity"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Service struct {
	client   *PlaywrightClient
	repo     *Repository
	headless bool
}

func NewService(headless bool, repo *Repository) (*Service, error) {
	client, err := NewPlaywrightClient(headless)
	if err != nil {
		return nil, err
	}
	return &Service{
		repo:     repo,
		client:   client,
		headless: headless,
	}, nil
}

func (s *Service) Close() {
	s.client.Close()
}

func GetDefaultKeyword() []string {
	return []string{
		"Backend",
		"Frontend",
		"Fullstack",
		"Software Engineer",
		"Web Developer",
	}
}

// fungsi internal untuk dipakai di fungsi utama yaitu scrap and save. fungsi ini digunakan untuk melakukan scraping ke beberapa keyword
func (s *Service) scrapeMultipleKeyword(ctx context.Context, keywords []string, limitPerKeyword int) ([]RawJob, error) {
	//membuat channel untuk mengumpulkan hasil dari setiap go routine
	results := make(chan []RawJob, len(keywords))
	errors := make(chan error, len(keywords))

	for _, keyword := range keywords {
		go s.scrapeKeyword(ctx, keyword, limitPerKeyword, results, errors)
	}

	//membuat wadah untuk nantinya dimasukkan data hasil scraping pada channel
	//membutuhkan wadah, karena setiap channel harus ditutup ketika sudah tidak dipakai. jika langsung mereturn channel, itu akan menyebabkan memory leak. alias kebocoran terhadap memory laptop
	var allJobs []RawJob
	var errs []error

	//melakukan looping. looping ini akan menunggu hasil dari worker, baru melanjutkan ke iterasi selanjutnyaa. dan memasukkannya kedalam wadah yang sudah dibuat
	for i := 0; i < len(keywords); i++ {
		select {
		case jobs := <-results:
			allJobs = append(allJobs, jobs...)
		case err := <-errors:
			errs = append(errs, err)
		case <-ctx.Done():
			return allJobs, ctx.Err()
		}
	}

	//jika ada error kembalikan alljobs dan error yang terjadi
	if len(errs) > 0 {
		return allJobs, fmt.Errorf("Beberapa keyword gagal: %v", errs)
	}

	return allJobs, nil
}

// fungsi ini dipakai untuk menghandle detail job scraping. dan membuat fungsinya menjadi private karena hanya dipakai di dalam fungsi internal service
func (s *Service) fetchJobDetail(ctx context.Context, rawJobs []RawJob) ([]RawJob, error) {
	//jika rawjobs kosong berhenti
	if len(rawJobs) == 0 {
		return []RawJob{}, nil
	}

	//membuat channel untuk antrian
	jobChan := make(chan RawJob, len(rawJobs))

	//membuat channel untuk hasil scrapping
	resultChan := make(chan RawJob, len(rawJobs))

	//membuat channel untuk error yang terjadi saat proses scraping
	errorChan := make(chan error, len(rawJobs))

	//memasukkan list rawjob kedalam antrian untuk nantinyua dieksekusi oleh go routines
	for _, job := range rawJobs {
		jobChan <- job
	}
	//menutup channel job, karna channel job sudah beres terisi
	close(jobChan)

	var wg sync.WaitGroup

	//melakukan looping sebanyak jumlah worker yang akan digunakan
	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(workerId int) {
			defer wg.Done()

			for job := range jobChan {
				select {
				case <-ctx.Done():
					errorChan <- ctx.Err()
					return
				default:
				}

				//melakukan scrap terhadap detal halaman, memanggil fungsi scrap detail
				detail, err := s.scrapeDetail(ctx, job.Url)
				if err != nil {
					log.Printf("gagal mengambil detail job %s:%v", job.Url, err)
					continue
				}

				//update job dengan detail yang sudah didapatkan
				job.Description = detail.Description
				job.Skills = detail.Skills

				//memasukkan hasil detail kedalam job channel
				resultChan <- job

				//menerapkan random delay untuk mengecoh keamanan glints
				RandomDelay(2000, 5000)
			}
		}(i)
	}

	go func() {
		wg.Wait()

		close(resultChan)
		close(errorChan)
	}()

	//inisialisasi untuk menyimpan semua detail job dari hasil resultchan
	var completedJob []RawJob
	//inisialisasi untuk menyimpan error dari hasil error chan
	var errors []error

	for result := range resultChan {
		completedJob = append(completedJob, result)
	}

	for err := range errorChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return completedJob, fmt.Errorf("beberapa detail gagal : %w", errors)
	}

	return completedJob, nil
}

// fungsi ini untuk melakukan 3 tugas
// 1. melakukan scrapping overview job
// 2. melakukan scrapping terhadap detail job di setiap link dari hasil overview job
// 3. melakukan save ke database
// fungsi ini public karna akan dipakai di handler
func (s *Service) ScrapeAndSave(ctx context.Context, userId string, targetPerKeyword int) (*ScraperResponse, error) {
	startTime := time.Now()

	//memanggil fungsi yang sudah didefinisikan, fungsi ini mengembalikan default keyword
	keywords := GetDefaultKeyword()

	// ini adalah phase pertama pada scraping.
	// pada fase ini mengambil overview awal pada job
	rawJob, err := s.scrapeMultipleKeyword(ctx, keywords, targetPerKeyword)
	fmt.Println(err)
	if err != nil {
		return nil, fmt.Errorf("Phase 1 failed : %w", err)
	}

	// ini adalah phase kedua pada scraping
	// pada fase ini mulai menelusuri ke detail page. untuk mendapatkan detail dair setiap jobnya
	completeJobs, err := s.fetchJobDetail(ctx, rawJob)
	if err != nil {
		return nil, fmt.Errorf("Phase 2 failed : %w", err)
	}

	//membuat wadah untuk memasukkan rawjob kedalam bentuk entity.job, untuk dimasukan kedalam database
	var jobs []entity.Job
	for _, raw := range completeJobs {
		//melakukan pengecekan, apakah lokasi mengandung kalimat remote
		locLower := strings.ToLower(raw.Location)
		isRemote := strings.Contains(locLower, "remote") ||
			strings.Contains(locLower, "wfh") ||
			strings.Contains(locLower, "work from home")

		//jika lokasi mengandung bandung, masukkan kota bandung, jika tidak mengandung kota bandung. kosongkan
		city := ""
		if strings.Contains(locLower, "bandung") {
			city = "Bandung"
		}

		skills := pq.StringArray(raw.Skills)

		id := uuid.New().String()

		//memasukkan data rawjob kedalam atribut entity job
		job := entity.Job{
			Id:          id,
			Title:       raw.Title,
			Company:     raw.Company,
			Description: raw.Description,
			Location:    raw.Location,
			IsRemote:    isRemote,
			Url:         raw.Url,
			City:        city,
			Source:      "glints",
			Skills:      skills,
			UserId:      userId,
		}

		//memasukkan job kedalam array, untuk nanti dikumpulkan dan dimasukkan
		jobs = append(jobs, job)
	}

	//memasukkan data jobs kedalam database
	inserted, duplicated, err := s.repo.SaveJobs(ctx, jobs)
	if err != nil {
		return nil, fmt.Errorf("save to database failed :%w", err)
	}

	executionTime := int(time.Since(startTime).Seconds())

	return &ScraperResponse{
		Message:          "Scraping completed",
		TotalFound:       len(completeJobs),
		TotalInserted:    inserted,
		TotalDuplicated:  duplicated,
		ExecutionTimeSec: executionTime,
	}, nil
}
