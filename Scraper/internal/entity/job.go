package entity

import (
	"time"

	"github.com/lib/pq"
)

type Job struct {
	Id          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string
	Company     string `gorm:"not null"`
	Description string `gorm:"not null"`
	Url         string `gorm:"uniqueIndex;not null"`
	Location    string `gorm:"not null"`
	City        string
	Source      string
	Skills      pq.StringArray `gorm:"type:text[]"`
	Salary      string
	ScrapedAt   time.Time `gorm:"autoCreateTime"`
}

func (Job) TableName() string {
	return "jobs"
}
