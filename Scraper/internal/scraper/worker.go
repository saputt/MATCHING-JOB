package scraper

import (
	"context"
	"fmt"
	"strings"

	"github.com/playwright-community/playwright-go"
)

// fungsi yang akan dijalan oleh worker
// dimana fungsi ini untuk melakukan scraping terhadap keyword yang diberikan
// hasil dari scrap ini adalah overview awal untuk setiap job
func (s *Service) scrapeKeyword(
	ctx context.Context,
	keyword string,
	target int,
	result chan<- []RawJob,
	errors chan<- error,
) {
	//membuat page baru
	page, err := s.client.NewPage()
	if err != nil {
		errors <- fmt.Errorf("keyword %s: failed to create page: %w", keyword, err)
		return
	}
	defer s.client.ClosePage(page)

	//melakukan blokir terhadap resource yang tidak perlu.
	//tidak memblokir css karena bisa menyebabkan kecurigaan
	// err = page.Route("**/*", func(route playwright.Route) {
	// 	resourceType := route.Request().ResourceType()

	// 	if resourceType == "image" || resourceType == "font" {
	// 		route.Abort()
	// 	} else {
	// 		route.Continue()
	// 	}
	// })

	if err != nil {
		errors <- fmt.Errorf("keyword %s: failed to set route: %w", keyword, err)
		return
	}

	// Ganti blok pembuatan URL kamu jadi sebersih ini:
	escapedKeyword := strings.ReplaceAll(keyword, " ", "+")
	url := fmt.Sprintf(
		"https://glints.com/id/opportunities/jobs/explore?country=ID&locationId=86a3dc56-1bd7-4cd3-8225-d3e4b976e552&locationName=Bandung,Jawa+Barat&lowestLocationLevel=3&keyword=%s",
		escapedKeyword,
	)

	//melakukan navigasi ke url, pada page yang sudah dibuat
	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateCommit,
		Timeout:   playwright.Float(60000),
	})

	if err != nil {
		errors <- fmt.Errorf("keyword %s: failed to navigate: %w", keyword, err)
		return
	}

	//menunggu sampai jaringan idle, atau sampai tidak ada request yang masuk
	//ini lebih baik daripada menggunakan time.sleep karena lebih adaptif dan tidak kaku
	// err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
	// 	State: playwright.LoadStateDomcontentloaded,
	// })

	// if err != nil {
	// 	errors <- fmt.Errorf("keyword %s: wait for load failed: %w", keyword, err)
	// 	return
	// }

	//menggunakan randow delay agar meniru perilaku manusia
	//agar menghindari deteksi bot oleh web glints
	RandomDelay(2000, 4000)

	cardLocator := page.Locator("div[class*='JobCardsc__JobcardContainer']")

	err = cardLocator.First().WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(60000),
	})

	if err != nil {
		errors <- fmt.Errorf("keyword %s: no job card found: %w", keyword, err)
		return
	}

	var jobs []RawJob
	seenUrl := make(map[string]bool)
	scrollCount := 0
	maxScroll := 15
	noNewJobCount := 0
	maxNoNewJob := 3

	//melakukan perulangan sampai kondisi job target terpenuhi atau maximal scroll telah terpenuhi
	//melakukan loop scraping dengan menerapkan human like scrolling
	for len(jobs) < target && scrollCount < maxScroll {
		//scroll kebawah
		_, err := page.Evaluate("window.scrollBy(0, window.innerHeight)")

		if err != nil {
			errors <- fmt.Errorf("keyword %s: scroll error: %w", keyword, err)
			return
		}

		RandomDelay(1000, 2000)
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

			raw := RawJob{}

			titleLocator := card.Locator("h2[class*='CompactOpportunityCardsc__JobTitle'] a")

			count, _ := titleLocator.Count()
			if count > 0 {
				titleText, err := titleLocator.First().TextContent()
				if err == nil {
					raw.Title = strings.TrimSpace(titleText)
				}
			}

			if !isItJob(raw.Title) {
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

			if isSusCompany(raw.Company) {
				continue
			}

			locationLocator := card.Locator("div[class*='CardJobLocation__LocationWrapper'] span[title]")

			count, _ = locationLocator.Count()
			if count > 0 {
				locationText, err := locationLocator.First().TextContent()
				if err == nil {
					raw.Location = locationText
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

		RandomDelay(500, 1500)
	}

	result <- jobs
}
