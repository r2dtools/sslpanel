package db

import (
	"backend/config"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func GetDB(config *config.Config) (*gorm.DB, error) {
	var err error

	if db == nil {
		db, err = gorm.Open(config.DbType, config.DbDsn)

		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
