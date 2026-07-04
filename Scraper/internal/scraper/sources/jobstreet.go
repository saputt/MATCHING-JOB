package sources

import (
	"context"
	"fmt"
	"job-matching-scraper/internal/client"
	"job-matching-scraper/internal/model"
	"job-matching-scraper/internal/utils"
	"log"
	"strings"
	"sync"

	"github.com/playwright-community/playwright-go"
)

type JobstreetSouce struct {
	client *client.PlaywrightClient
}

func NewJobStreetSource(client *client.PlaywrightClient) *JobstreetSouce {
	return &JobstreetSouce{
		client: client,
	}
}

func (s *JobstreetSouce) GetName() string {
	return "jobstreet"
}

func (s *JobstreetSouce) Scrape(ctx context.Context, keywords []string, limit int, existingJobMap map[string]bool) ([]model.RawJob, error) {
	result := make(chan []model.RawJob, len(keywords))
	errors := make(chan error, len(keywords))

	maxConcurent := 7
	semaphore := make(chan struct{}, maxConcurent)

	for _, keyword := range keywords {
		semaphore <- struct{}{}
		go func(kw string) {
			defer func() { <-semaphore }()
			jobs, err := s.scrapeKeyword(ctx, kw, limit, existingJobMap)
			if err != nil {
				errors <- fmt.Errorf("glints %s : %w", kw, err)
			}
			result <- jobs
		}(keyword)
	}

	var allRawJob []model.RawJob

	for i := 0; i < len(keywords); i++ {
		select {
		case job := <-result:
			allRawJob = append(allRawJob, job...)
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
	for i := 0; i < 15; i++ {
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

				detail, err := s.ScrapeDetail(job.Url)
				if err != nil {
					continue
				}

				job.Description = detail.Description
				job.Salary = detail.Salary

				resultChan <- job

				utils.RandomDelay(1500, 3000)
			}

		}(i)
		utils.RandomDelay(1500, 3000)
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

func (s *JobstreetSouce) scrapeKeyword(ctx context.Context, keyword string, limit int, existingJobMap map[string]bool) ([]model.RawJob, error) {
	page, err := s.client.NewPage()
	if err != nil {
		return nil, err
	}
	defer s.client.ClosePage(page)

	var jobs []model.RawJob
	noNewJobCount := 0
	maxNoNewJob := 3
	currentPage := 1
	maxPage := 20

	url := fmt.Sprintf(
		"https://id.jobstreet.com/%s-jobs/in-Bandung-West-Java?page=%d",
		strings.ReplaceAll(strings.ToLower(keyword), " ", "-"),
		currentPage,
	)

	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateCommit,
	})

	if err != nil {
		return nil, err
	}

	utils.RandomDelay(2000, 4000)

	cardLocator := page.Locator("article[data-automation='normalJob']")

	err = cardLocator.First().WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(30000),
	})

	if err != nil {
		log.Println("failed to load jobstreet page, block by cloudflare")
		return nil, err
	}

	for len(jobs) < limit && currentPage < maxPage {
		select {
		case <-ctx.Done():
			return jobs, ctx.Err()
		default:
		}

		utils.RandomDelay(1000, 2000)

		cards, err := cardLocator.All()
		if err != nil {
			continue
		}

		newJobCount := 0

		for _, card := range cards {
			if len(jobs) >= limit {
				break
			}

			raw := model.RawJob{}

			titleLoc := card.Locator("[data-automation='jobTitle']")
			count, _ := titleLoc.Count()
			if count > 0 {
				titleText, err := titleLoc.First().TextContent()
				log.Println("jobstreet scrape. keyword : " + keyword + ", title : " + titleText)
				if err == nil {
					raw.Title = strings.TrimSpace(titleText)
				}
			}

			if !utils.IsItJob(raw.Title) {
				continue
			}

			companyLoc := card.Locator("[data-automation='jobCompany']")
			count, _ = companyLoc.Count()
			if count > 0 {
				companyText, err := companyLoc.First().TextContent()
				if err == nil {
					raw.Company = strings.TrimSpace(companyText)
				}
			}

			key := raw.Title + "|" + raw.Company

			if existingJobMap[key] {
				continue
			}

			if utils.IsSusCompany(raw.Company) {
				continue
			}

			locationLoc := card.Locator("[data-automation='jobCardLocation']")
			count, _ = locationLoc.Count()
			if count > 0 {
				locationText, err := locationLoc.First().TextContent()
				if err == nil {
					raw.Location = strings.TrimSpace(locationText)
				}
			}

			linkLoc := card.Locator("a[data-automation='jobTitle']")
			count, _ = linkLoc.Count()
			if count > 0 {
				linkHref, _ := linkLoc.First().GetAttribute("href")
				if linkHref != "" {
					raw.Url = "https://id.jobstreet.com" + linkHref
				}
			}

			raw.Source = "jobstreet"

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

			log.Printf("scraping jobstreet success. title : %s, company : %s", raw.Title, raw.Company)
		}

		if newJobCount == 0 {
			noNewJobCount++
			if noNewJobCount >= maxNoNewJob {
				break
			}
		}

		utils.RandomDelay(1500, 3000)

		currentPage++
	}

	utils.RandomDelay(3000, 6000)
	return jobs, nil
}

// fungsi ini merupakan tahapan kedua dari scraping yang berfungsi untuk melakukan scraping pada detail job yang sudah discraping pada tahap satu
func (s *JobstreetSouce) ScrapeDetail(
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

	salaryLocator := page.Locator(`span[data-automation="job-detail-salary"]`)
	count, _ := salaryLocator.Count()
	if count > 0 {
		salaryText, err := salaryLocator.First().TextContent()
		if err == nil {
			detail.Salary = strings.TrimSpace(salaryText)
		}
	}

	descLocator := page.Locator("div[data-automation='jobAdDetails']")

	err = descLocator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(5000),
	})

	if err == nil {
		descContent, err := descLocator.InnerHTML()
		if err == nil {
			detail.Description = utils.CleanJobstreetDesc(descContent)
		}
	} else {
		detail.Description = ""
	}

	utils.RandomDelay(1500, 2000)

	return detail, nil
}
