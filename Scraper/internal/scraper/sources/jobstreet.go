package sources

import (
	"context"
	"fmt"
	"job-matching-scraper/internal/client"
	"job-matching-scraper/internal/model"
	"job-matching-scraper/internal/utils"
	"strings"

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

func (s *JobstreetSouce) Scrape(ctx context.Context, keywords []string, limit int) ([]model.RawJob, error) {
	result := make(chan []model.RawJob, len(keywords))
	errors := make(chan error, len(keywords))

	maxConcurent := 5
	semaphore := make(chan struct{}, maxConcurent)

	for _, keyword := range keywords {
		semaphore <- struct{}{}
		go func(kw string) {
			defer func() { <-semaphore }()
			jobs, err := s.scrapeKeyword(ctx, kw, limit)
			if err != nil {
				errors <- fmt.Errorf("glints %s : %w", kw, err)
			}
			result <- jobs
		}(keyword)
	}

	var allJobs []model.RawJob

	for i := 0; i < len(keywords); i++ {
		select {
		case job := <-result:
			allJobs = append(allJobs, job...)
		case err := <-errors:
			fmt.Printf("Glints error: %v\n", err)
		case <-ctx.Done():
			return allJobs, ctx.Err()
		}
	}

	return allJobs, nil
}

func (s *JobstreetSouce) scrapeKeyword(ctx context.Context, keyword string, limit int) ([]model.RawJob, error) {
	page, err := s.client.NewPage()
	if err != nil {
		return nil, err
	}
	defer s.client.ClosePage(page)

	var jobs []model.RawJob
	seenUrl := make(map[string]bool)
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
		Timeout:   playwright.Float(30000),
	})

	if err != nil {
		return nil, err
	}

	utils.RandomDelay(2000, 4000)

	cardLocator := page.Locator("article[data-automation='normalJob']")

	err = cardLocator.First().WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(30000),
	})

	for len(jobs) < limit && currentPage < maxPage {
		select {
		case <-ctx.Done():
			return jobs, ctx.Err()
		default:
		}

		utils.RandomDelay(1500, 3000)

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

			descriptionLoc := card.Locator("div[data-automation='jobAdDetails'] div")
			count, _ = descriptionLoc.Count()
			if count > 0 {
				descriptionText, err := descriptionLoc.InnerHTML()
				if err == nil {
					raw.Description = strings.TrimSpace(descriptionText)
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

		currentPage++
	}

	utils.RandomDelay(3000, 6000)
	return jobs, nil
}
