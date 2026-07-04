package database

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(databaseUrl string) (*gorm.DB, error) {
	if databaseUrl == "" {
		return nil, errors.New("Database url cannot be empty")
	}

	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
