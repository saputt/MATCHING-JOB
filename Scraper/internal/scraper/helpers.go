package scraper

import (
	"math/rand"
	"strings"
	"time"
)

// fungsi ini untuk melaukan random delay, yang digunakan untuk mengelabui sistem
func RandomDelay(minMs, maxMs int) {
	if minMs < 0 {
		minMs = 2000
	}
	if maxMs < 0 {
		maxMs = 5000
	}

	delay := minMs + rand.Intn(maxMs-minMs)
	jitter := delay / 10
	finalDelay := delay + rand.Intn(jitter) - (jitter / 2)

	time.Sleep(time.Duration(finalDelay) * time.Millisecond)
}

func GetStealthArgs() []string {
	return []string{
		// Menonaktifkan fitur yang menandai browser sebagai automation
		// Tanpa ini, navigator.webdriver akan bernilai true
		"--disable-blink-features=AutomationControlled",

		// Menonaktifkan isolasi origin (seperti browser normal)
		"--disable-features=IsolateOrigins,site-per-process",

		// Menonaktifkan web security (tidak masalah untuk scraping)
		"--disable-web-security",

		// Menonaktifkan fitur yang bisa mendeteksi headless
		"--disable-features=BlockInsecurePrivateNetworkRequests",

		// Menonaktifkan sync (browser normal juga bisa)
		"--disable-sync",

		// Menonaktifkan default apps (biar lebih clean)
		"--disable-default-apps",

		// Menonaktifkan extensions (browser normal kadang juga)
		"--disable-extensions",

		"--no-sandbox",
	}
}

// fungsi ini untuk mengecek apakah company, merupakan company yang mencurigakan dan melenceng dari IT
func isSusCompany(company string) bool {
	lowerCompany := strings.ToLower(company)

	suspicious := []string{
		"recruitment", "agency", "outsourcing", "consulting",
		"manpower", "staffing", "headhunter", "talent",
	}

	for _, s := range suspicious {
		if strings.Contains(lowerCompany, s) {
			return true
		}
	}

	return false
}

func GetRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
	return userAgents[rand.Intn(len(userAgents))]
}

// fungsi ini untuk mengecek, apakah kata kunci job yang didapat dari scrap merupakan job IT
func isItJob(title string) bool {
	lowerTitle := strings.ToLower(title)

	itKeyword := []string{
		"software", "developer", "engineer", "programmer",
		"backend", "frontend", "fullstack", "full-stack",
		"web", "mobile", "android", "ios", "devops",
		"data", "machine learning", "ai", "system analyst",
		"it", "information technology", "cloud", "security",
	}

	for _, keyword := range itKeyword {
		if strings.Contains(lowerTitle, keyword) {
			return true
		}
	}

	return false
}

// fungis ini untuk mengecek kemungkinan dari setiap role. misal hasil scraping bisa beragam, misal go, golang dll. fungsi ini untuk mengecek beberapa kemungkinan
func normalizeSkill(skill string) string {
	lowerSkill := strings.ToLower(strings.TrimSpace(skill))

	synonyms := map[string][]string{
		"golang":     {"go", "go language", "golang", "go-lang"},
		"javascript": {"javascript", "js", "ecmascript"},
		"typescript": {"typescript", "ts"},
		"react":      {"react", "react.js", "reactjs", "react-js"},
		"nodejs":     {"node", "node.js", "nodejs", "express"},
		"python":     {"python", "py"},
		"java":       {"java", "openjdk"},
		"docker":     {"docker", "container", "dockerfile"},
		"postgresql": {"postgres", "postgresql", "pgsql"},
		"mysql":      {"mysql", "mariadb"},
		"redis":      {"redis", "cache"},
		"git":        {"git", "github", "gitlab", "version control"},
	}

	for canonical, variants := range synonyms {
		for _, variant := range variants {
			if lowerSkill == variant || strings.Contains(lowerSkill, variant) {
				return canonical
			}
		}
	}

	return lowerSkill
}
