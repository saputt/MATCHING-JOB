package entity

import "time"

type User struct {
	Id                 string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Username           string    `gorm:"column:username;type:varchar(255);not null"`
	Email              string    `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	Password           string    `gorm:"column:password;type:varchar(255);not null"`
	LocationPreference string    `gorm:"column:location_preference;default:ALL"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Relasi Has Many: 1 User bisa punya banyak Job
	Jobs []Job `gorm:"foreignKey:UserId;references:Id"`
}

func (User) TableName() string {
	return "users"
}
