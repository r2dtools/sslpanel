package db

import (
	"backend/config"
	"github.com/jinzhu/gorm"
	// Initialize postgresql dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

// Init database connection
func Init() error {
	var err error
	config := config.GetConfig()
	db, err = gorm.Open("postgres", config.DatabaseURI)

	if err != nil {
		return err
	}

	return nil
}

// GetDB return db connection pool
func GetDB() *gorm.DB {
	return db
}
