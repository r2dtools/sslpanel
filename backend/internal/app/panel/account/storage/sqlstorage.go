package storage

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type sqlStorage struct {
	db *gorm.DB
}

func (s *sqlStorage) FindById(id uint) (*Account, error) {
	var acc Account
	err := s.db.First(&acc, id).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not find account with ID %d: %w", id, err)
	}

	return &acc, nil
}

func (s *sqlStorage) Save(a *Account) error {
	if s.db.NewRecord(a) {
		return s.db.Create(a).Error
	}

	return s.db.Save(a).Error
}

func NewAccountSqlStorage(db *gorm.DB) AccountStorage {
	return &sqlStorage{
		db: db,
	}
}
