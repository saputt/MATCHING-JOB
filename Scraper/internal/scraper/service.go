package scraper

import (
	"context"
	"fmt"
	"job-matching-scraper/internal/client"
	"job-matching-scraper/internal/entity"
	"job-matching-scraper/internal/model"
	"job-matching-scraper/internal/scraper/sources"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Service struct {
	client  *client.PlaywrightClient
	repo    *Repository
	sources []sources.JobSource
}

func NewService(headless bool, repo *Repository) (*Service, error) {
	client, err := client.NewPlaywrightClient(headless)
	if err != nil {
		return nil, err
	}

	allSources := []sources.JobSource{
		sources.NewGlintsSource(client),
		sources.NewJobStreetSource(client),
	}

	return &Service{
		repo:    repo,
		client:  client,
		sources: allSources,
	}, nil
}

func (s *Service) Close() {
	s.client.Close()
}

func GetDefaultKeyword() []string {
	return []string{
		// --- Role Utama (Punya Kamu) ---
		"Backend", "Frontend", "Fullstack", "Software Engineer",
		"Web Developer", "DevOps", "Data Engineer", "Mobile Developer",

		// --- Specific Seniority (Fresh Grad / Intern) ---
		"Junior Developer", "Junior Backend", "Junior Frontend",
		"Internship Software Engineer", "Magang Developer", "Entry Level Developer",
		"Associate Software Engineer",

		// --- Tech Stack Spesifik ---
		"React Developer", "Node.js Developer", "Golang Developer",
		"Laravel Developer", "Flutter Developer",

		// --- Perluasan Role Tech ---
		"QA Engineer", "Data Analyst", "IT Support", "Technical Writer",
	}
}

// fungsi ini untuk melakukan 3 tugas
// 1. melakukan scrapping overview job
// 2. melakukan scrapping terhadap detail job di setiap link dari hasil overview job
// 3. melakukan save ke database
// fungsi ini public karna akan dipakai di handler
func (s *Service) ScrapeAndSave(ctx context.Context, userId string, targetPerKeyword int) (*model.ScrapeResponse, error) {
	startTime := time.Now()

	//memanggil fungsi yang sudah didefinisikan, fungsi ini mengembalikan default keyword
	keywords := GetDefaultKeyword()

	results := make(chan []model.RawJob, len(keywords))
	errors := make(chan error, len(keywords))

	for _, source := range s.sources {
		go func(src sources.JobSource) {
			log.Printf("Starting scrape from source : %s", src.GetName())
			jobs, err := src.Scrape(ctx, keywords, targetPerKeyword)
			if err != nil {
				errors <- fmt.Errorf("%s: %w", src.GetName(), err)
				return
			}
			log.Printf("Source %s: %d jobs found", src.GetName(), len(jobs))
			results <- jobs
		}(source)
	}

	var allJobs []model.RawJob
	var errs []error

	for i := 0; i < len(s.sources); i++ {
		select {
		case jobs := <-results:
			allJobs = append(allJobs, jobs...)
		case err := <-errors:
			errs = append(errs, err)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if len(errs) > 0 && len(allJobs) == 0 {
		return nil, fmt.Errorf("all source failed: %v", errs)
	}

	//membuat wadah untuk memasukkan rawjob kedalam bentuk entity.job, untuk dimasukan kedalam database
	var jobs []entity.Job
	for _, raw := range allJobs {
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
			Source:      raw.Source,
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

	return &model.ScrapeResponse{
		Message:          "Scraping completed",
		TotalFound:       len(allJobs),
		TotalInserted:    inserted,
		TotalDuplicated:  duplicated,
		ExecutionTimeSec: executionTime,
	}, nil
}
