package db

import (
	"backend/config"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

var db *gorm.DB

func GetDB(config *config.Config) (*gorm.DB, error) {
	var err error

	if db == nil {
		db, err = gorm.Open(mysql.Open(config.DbDsn), &gorm.Config{})

		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
