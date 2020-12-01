package util

import (
	schema "github.com/baseli/fb_crm_audience/api/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
)

func NewDatabase() (*gorm.DB, error) {
	filepath, err := GetHomePath()
	if err != nil {
		return nil, err
	}

	filepath = path.Join(filepath, "database.db")
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&schema.Account{}, &schema.AdAccount{}, &schema.Task{})

	return db, nil
}

