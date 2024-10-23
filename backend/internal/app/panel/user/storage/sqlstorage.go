package storage

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type sqlStorage struct {
	db *gorm.DB
}

func (s *sqlStorage) FindById(id int) (*User, error) {
	var user User
	err := s.db.Preload("Account").First(&user, id).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not find user with ID %d: %w", id, err)
	}

	return &user, nil
}

func (s *sqlStorage) FindByEmail(email string) (*User, error) {
	var user User
	err := s.db.First(&user, "email = ?", email).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not find user with email %s: %w", email, err)
	}

	return &user, nil
}

func (s *sqlStorage) FindAll() ([]*User, error) {
	users := []*User{}

	if err := s.db.Preload("Account").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *sqlStorage) Save(a *User) error {
	if s.db.NewRecord(a) {
		return s.db.Create(a).Error
	}

	return s.db.Save(a).Error
}

func NewUserSqlStorage(db *gorm.DB) UserStorage {
	return &sqlStorage{
		db: db,
	}
}
