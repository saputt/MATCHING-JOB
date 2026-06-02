package sources

import (
	"context"
	"job-matching-scraper/internal/model"
)

type JobSource interface {
	GetName() string
	Scrape(ctx context.Context, keywords []string, limit int) ([]model.RawJob, error)
}
