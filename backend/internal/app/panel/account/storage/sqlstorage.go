package storage

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type sqlStorage struct {
	db *gorm.DB
}

func (s *sqlStorage) FindById(id uint) (*Account, error) {
	var acc Account
	err := s.db.First(&acc, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not find account with ID %d: %w", id, err)
	}

	return &acc, nil
}

func (s *sqlStorage) Save(a *Account) error {
	if a.ID == 0 {
		return s.db.Create(a).Error
	}

	return s.db.Save(a).Error
}

func NewAccountSqlStorage(db *gorm.DB) AccountStorage {
	return &sqlStorage{
		db: db,
	}
}
