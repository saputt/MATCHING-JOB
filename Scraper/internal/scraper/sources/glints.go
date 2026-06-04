package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"job-matching-scraper/internal/client"
	"job-matching-scraper/internal/model"
	"job-matching-scraper/internal/utils"
	"log"
	"strings"
	"sync"

	"github.com/playwright-community/playwright-go"
)

type GlintsSource struct {
	client *client.PlaywrightClient
}

func NewGlintsSource(client *client.PlaywrightClient) *GlintsSource {
	return &GlintsSource{
		client: client,
	}
}

func (s *GlintsSource) GetName() string {
	return "glints"
}

// fungsi ini adalah fungsi utama scraping terhadap glints
func (s *GlintsSource) Scrape(ctx context.Context, keywords []string, limit int) ([]model.RawJob, error) {
	// membuat channgel untuk menampung hasil dari setiap go routine
	results := make(chan []model.RawJob, len(keywords))
	// membaut channel unutk menampung error dari setiap go routine berdasarkan keyword yang di scrape
	errors := make(chan error, len(keywords))

	maxConcurent := 5
	semaphore := make(chan struct{}, maxConcurent)

	// melakukan looping berdasarkan jumlah keyword, laluu menjalankan go routines untuk melakukan scraping
	for _, keyword := range keywords {
		semaphore <- struct{}{}
		go func(kw string) {
			defer func() { <-semaphore }()
			jobs, err := s.scrapeKeyword(ctx, kw, limit)
			if err != nil {
				errors <- fmt.Errorf("glints %s : %w", kw, err)
			}
			results <- jobs
		}(keyword)
	}

	var allRawJob []model.RawJob

	// melakukan looping berdasarkan jumlah keyword, untuk mengumpulkan hasil dari masing masing worker dan memasukkannya kedalam allRawJob
	for i := 0; i < len(keywords); i++ {
		select {
		case jobs := <-results:
			allRawJob = append(allRawJob, jobs...)
		case err := <-errors:
			fmt.Printf("Glints error: %v\n", err)
		case <-ctx.Done():
			return allRawJob, ctx.Err()
		}
	}

	// membuat channel sebanyak panjang data yang berhasil di scraping sebelumnya
	jobChan := make(chan model.RawJob, len(allRawJob))
	resultChan := make(chan model.RawJob, len(allRawJob))
	errorsChan := make(chan error, len(allRawJob))

	// melakukan looping sebanyak data yang telah didapatkan lalu memasukkannya ke jobChan, yang akan digunakan sebagai antrian untuk goroutines
	for _, job := range allRawJob {
		jobChan <- job
	}
	close(jobChan)

	var wg sync.WaitGroup

	// melakukan looping sebanyak jumlah worker yang diinginkan, lalu mengerjakan scraping detail
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			for job := range jobChan {
				select {
				case <-ctx.Done():
					errorsChan <- ctx.Err()
					return
				default:
				}

				log.Printf("scraping detail jobstreet. title : %s", job.Title)

				detail, err := s.scrapeDetail(job.Url)
				if err != nil {
					continue
				}

				job.Description = detail.Description
				job.Skills = detail.Skills
				job.Location = detail.Location

				resultChan <- job

				utils.RandomDelay(1500, 3000)
			}

		}(i)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorsChan)
	}()

	var completeJobs []model.RawJob

	for job := range resultChan {
		completeJobs = append(completeJobs, job)
	}

	for err := range errorsChan {
		fmt.Printf("Detail scraping error: %v\n", err)
	}

	return completeJobs, nil
}

// fungsi yang akan dijalan oleh worker
// dimana fungsi ini untuk melakukan scraping terhadap keyword yang diberikan
// hasil dari scrap ini adalah overview awal untuk setiap job
func (s *GlintsSource) scrapeKeyword(
	ctx context.Context,
	keyword string,
	target int,
) ([]model.RawJob, error) {
	//membuat page baru
	page, err := s.client.NewPage()
	if err != nil {
		return nil, err
	}
	defer s.client.ClosePage(page)

	// Ganti blok pembuatan URL kamu jadi sebersih ini:
	escapedKeyword := strings.ReplaceAll(keyword, " ", "+")
	url := fmt.Sprintf(
		"https://glints.com/id/opportunities/jobs/explore?country=ID&locationId=86a3dc56-1bd7-4cd3-8225-d3e4b976e552&locationName=Bandung,Jawa+Barat&lowestLocationLevel=3&keyword=%s",
		escapedKeyword,
	)

	//melakukan navigasi ke url, pada page yang sudah dibuat
	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateCommit,
		Timeout:   playwright.Float(30000),
	})

	if err != nil {
		return nil, err
	}

	//menggunakan randow delay agar meniru perilaku manusia
	//agar menghindari deteksi bot oleh web glints
	utils.RandomDelay(2000, 4000)

	cardLocator := page.Locator("div[class*='JobCardsc__JobcardContainer']")

	err = cardLocator.First().WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(30000),
	})

	if err != nil {
		log.Println("failed to load jobstreet page, block by cloudflare")
		return nil, err
	}

	var jobs []model.RawJob
	seenUrl := make(map[string]bool)
	scrollCount := 0
	maxScroll := 20
	noNewJobCount := 0
	maxNoNewJob := 3

	//melakukan perulangan sampai kondisi job target terpenuhi atau maximal scroll telah terpenuhi
	//melakukan loop scraping dengan menerapkan human like scrolling
	for len(jobs) < target && scrollCount < maxScroll {
		select {
		case <-ctx.Done():
			return jobs, ctx.Err()
		default:
		}
		//scroll kebawah
		_, err := page.Evaluate("window.scrollBy(0, window.innerHeight)")

		if err != nil {
			return nil, err
		}

		utils.RandomDelay(1000, 2000)
		scrollCount++

		//mengambil semua kartu pada halaman
		cards, err := cardLocator.All()
		if err != nil {
			continue
		}

		newJobCount := 0

		for _, card := range cards {
			if len(jobs) >= target {
				break
			}

			raw := model.RawJob{}

			titleLocator := card.Locator("h2[class*='CompactOpportunityCardsc__JobTitle'] a")

			count, _ := titleLocator.Count()
			if count > 0 {
				titleText, err := titleLocator.First().TextContent()
				log.Printf("scraping glints. keyword : %s, title : %s ", keyword, titleText)
				if err == nil {
					raw.Title = strings.TrimSpace(titleText)
				}
			}

			if !utils.IsItJob(raw.Title) {
				continue
			}

			companyLocator := card.Locator("a[class*='CompactOpportunityCardsc__CompanyLinkResolver']")
			count, _ = companyLocator.Count()
			if count > 0 {
				companyText, err := companyLocator.First().TextContent()
				if err == nil {
					raw.Company = companyText
				}
			}

			if utils.IsSusCompany(raw.Company) {
				continue
			}

			locationLocator := card.Locator("div[class*='CardJobLocation__LocationWrapper']")

			count, _ = locationLocator.Count()
			if count > 0 {
				locationText, err := locationLocator.First().TextContent()
				if err == nil {
					raw.Location = utils.ExtractCityFromAdress(locationText)
				}
			}

			linkLocator := card.Locator("a")

			count, _ = linkLocator.Count()
			if count > 0 {
				href, _ := linkLocator.First().GetAttribute("href")
				if href != "" {
					raw.Url = "https://glints.com" + href
				}
			}

			if raw.Url == "" || seenUrl[raw.Url] {
				continue
			}

			raw.Source = "glints"

			seenUrl[raw.Url] = true

			locLower := strings.ToLower(raw.Location)
			isBandung := strings.Contains(locLower, "bandung")
			isRemote := strings.Contains(locLower, "remote") ||
				strings.Contains(locLower, "wfh") ||
				strings.Contains(locLower, "work from home")
			isHybrid := strings.Contains(locLower, "hybrid")

			if !isBandung && !isRemote && !isHybrid {
				continue
			}

			jobs = append(jobs, raw)
			newJobCount++
		}

		if newJobCount == 0 {
			noNewJobCount++
			if noNewJobCount >= maxNoNewJob {
				break
			}
		}

		utils.RandomDelay(1500, 3000)
	}

	utils.RandomDelay(3000, 6000)
	return jobs, nil
}

// fungsi ini merupakan tahapan kedua dari scraping yang berfungsi untuk melakukan scraping pada detail job yang sudah discraping pada tahap satu
func (s *GlintsSource) scrapeDetail(
	url string,
) (*model.JobDetail, error) {
	//membuat page baru, setiap scraping detail memakai page baru
	page, err := s.client.NewPage()

	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer s.client.ClosePage(page)

	//melakukan blokir terhadap resource yang tidak perlu
	err = page.Route("**/*", func(route playwright.Route) {
		resourceType := route.Request().ResourceType()

		if resourceType == "image" || resourceType == "font" {
			route.Abort()
		} else {
			route.Continue()
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to set route: %w", err)
	}

	page.AddInitScript(playwright.Script{
		Content: playwright.String(`
        Object.defineProperty(navigator, 'webdriver', { get: () => undefined });
        window.navigator.chrome = { runtime: {} };
        Object.defineProperty(navigator, 'languages', { get: () => ['en-US', 'en'] });
    `),
	})

	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(30000),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})

	if err != nil {
		return nil, fmt.Errorf("wait for load is failed: %w", err)
	}

	utils.RandomDelay(1000, 2000)

	detail := &model.JobDetail{}

	descLocator := page.Locator("div.JobDescriptionsc__DescriptionContainer div.DraftjsReadersc__ContentContainer")

	count, _ := descLocator.Count()
	if count > 0 {
		desc, err := descLocator.First().TextContent()
		if err == nil {
			detail.Description = strings.TrimSpace(desc)
		}
	}

	if detail.Description == "" {
		altDescLocator := page.Locator("div[class*='JobDescription'] div[class*='DraftjsReader']")
		count, _ := altDescLocator.Count()
		if count > 0 {
			desc, err := altDescLocator.First().TextContent()
			if err == nil {
				detail.Description = strings.TrimSpace(desc)
			}
		}
	}

	skillLocator := page.Locator("div.Skillssc__TagContainer div.TagStyle-sc-r1wv7a-4 span div")

	count, _ = skillLocator.Count()

	if count > 0 {
		allSkills, err := skillLocator.All()
		if err == nil {
			for _, skillEl := range allSkills {
				skillText, err := skillEl.TextContent()
				if err == nil {
					normalized := utils.NormalizeSkill(strings.TrimSpace(skillText))
					if normalized != "" {
						detail.Skills = append(detail.Skills, normalized)
					}
				}
			}
		}
	}

	//fallback jika skill tidak ditemukan
	if len(detail.Skills) == 0 {
		jsonLdLocator := page.Locator("script[type='application/ld+json']")
		count, _ := jsonLdLocator.Count()
		if count > 0 {
			for i := 0; i < count; i++ {
				script, err := jsonLdLocator.Nth(i).TextContent()

				if err == nil && strings.Contains(script, "JobPosting") {
					var jobPosting map[string]interface{}
					err = json.Unmarshal([]byte(script), &jobPosting)
					if err == nil {
						skillval, ok := jobPosting["skills"]
						if ok {
							skillStr, ok := skillval.(string)
							if ok {
								for _, s := range strings.Split(skillStr, ",") {
									normalize := utils.NormalizeSkill(strings.TrimSpace(s))
									if normalize != "" {
										detail.Skills = append(detail.Skills, normalize)
									}
								}
							}
						}
					}
					break
				}
			}
		}
	}

	locationLocator := page.Locator("div[class*='AboutCompanySectionsc__AddressWrapper'] p")
	count, _ = locationLocator.Count()
	if count > 0 {
		locationText, err := locationLocator.First().TextContent()
		if err == nil {
			detail.Location = strings.TrimSpace(locationText)
		}
	}

	utils.RandomDelay(1500, 2000)

	return detail, nil
}
