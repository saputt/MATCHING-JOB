package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/playwright-community/playwright-go"
)

// fungsi ini merupakan tahapan kedua dari scraping yang berfungsi untuk melakukan scraping pada detail job yang sudah discraping pada tahap satu
func (s *Service) scrapeDetail(
	ctx context.Context,
	url string,
) (*JobDetail, error) {
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

	_, err = page.Goto(url)

	if err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})

	if err != nil {
		return nil, fmt.Errorf("wait for load is failed: %w", err)
	}

	RandomDelay(1000, 2000)

	detail := &JobDetail{}

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
					normalized := normalizeSkill(strings.TrimSpace(skillText))
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
									normalize := normalizeSkill(strings.TrimSpace(s))
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

	companyLocator := page.Locator("div.TopFoldExperimentsc__JobOverViewCompanyName a")
	count, _ = companyLocator.Count()
	if count > 0 {
		companyText, err := companyLocator.First().TextContent()
		if err == nil {
			detail.Company = strings.TrimSpace(companyText)
		}
	}

	locationLocator := page.Locator("div.CardJobLocation__LocationWrapper span[title]")
	count, _ = locationLocator.Count()
	if count > 0 {
		locationText, err := locationLocator.First().TextContent()
		if err == nil {
			detail.Location = strings.TrimSpace(locationText)
		}
	}

	return detail, nil
}
