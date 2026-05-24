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
	IsRemote    bool   `gorm:"default:false"`
	City        string
	Source      string         `gorm:"default:glints"`
	Skills      pq.StringArray `gorm:"type:text[]"`
	ScrapedAt   time.Time      `gorm:"autoCreateTime"`
	UserId      string         `gorm:"column:user_id;type:uuid;not null"`
	User        *User          `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
}

func (Job) TableName() string {
	return "jobs"
}
