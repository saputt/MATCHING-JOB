package scraper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/playwright-community/playwright-go"
)

type PlaywrightClient struct {
	pw       *playwright.Playwright
	browser  playwright.Browser
	headless bool
}

func NewPlaywrightClient(headless bool) (*PlaywrightClient, error) {
	//start playwright, dan menginisialisasi browser driver
	pw, err := playwright.Run()

	if err != nil {
		return nil, err
	}

	_, err = os.Stat("state.json")

	if os.IsNotExist(err) {
		pw.Stop()
		return nil, fmt.Errorf("state.json not found. Please run LoginManual() first")
	}

	//konfigurasu launch options untuk browser
	launchOps := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	}

	//jika headless = false, tambahkan stealarg untuk menyamarkan fingerptint headless browser
	if !headless {
		launchOps.Args = []string{
			"--disable-blink-features=AutomationControlled",
		}
	}

	//launch chromunium
	browser, err := pw.Chromium.Launch(launchOps)

	//jika gagal tutup playwright
	if err != nil {
		pw.Stop()
		return nil, err
	}

	return &PlaywrightClient{
		pw:       pw,
		browser:  browser,
		headless: headless,
	}, nil
}

// fungsi untuk membuat page baru pada browser dengan dengan user agent random
// setiap page memiliki context yang berbeda
// hal ini penting untuk menghindari deteksi browser yang sama secara terus menerus
func (c *PlaywrightClient) NewPage() (playwright.Page, error) {
	stateBytes, err := os.ReadFile("state.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read state.json : %w", err)
	}

	var storageStateObj playwright.OptionalStorageState

	err = json.Unmarshal(stateBytes, &storageStateObj)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state.json into OptionalStorageState: %w", err)
	}

	//membuat context baru dengan user agent yang berbeda
	//context beriskan environtment seperti cookie, storage dll
	context, err := c.browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent:    playwright.String(GetRandomUserAgent()),
		StorageState: &storageStateObj,
	})

	if err != nil {
		return nil, err
	}

	//buka page baru dari context yang sudah dibuat
	page, err := context.NewPage()

	//melakukan cleanup jika gagal
	if err != nil {
		context.Close()
		return nil, err
	}

	return page, nil
}

func (c *PlaywrightClient) ClosePage(page playwright.Page) {
	if page != nil {
		ctx := page.Context()

		err := page.Close()
		if err != nil {
			fmt.Printf("failed to close page cleanly: %v\n", err)
		}

		if ctx != nil {
			ctx.Close()
		}
	}
}

func LoginManual() error {
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("failed to start playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})

	if err != nil {
		return fmt.Errorf("failed to launch browser : %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("failed to create page : %w", err)
	}

	_, err = page.Goto("https://glints.com/id")
	if err != nil {
		return fmt.Errorf("failed to direct to glints web : %w", err)
	}

	fmt.Scanln()

	storage, err := page.Context().StorageState()
	if err != nil {
		return fmt.Errorf("failed to get storage state : %w", err)
	}

	data, err := json.MarshalIndent(storage, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal state : %w", err)
	}

	err = os.WriteFile("state.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state.json : %w", err)
	}

	return nil
}

func (c *PlaywrightClient) Close() {
	if c.browser != nil {
		c.browser.Close()
	}
	if c.pw != nil {
		c.pw.Stop()
	}
}
