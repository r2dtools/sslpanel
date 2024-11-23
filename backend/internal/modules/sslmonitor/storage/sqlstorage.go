package storage

import (
	"encoding/json"
	"errors"

	"github.com/r2dtools/agentintegration"
	"gorm.io/gorm"
)

type sqlStorage struct {
	db *gorm.DB
}

func (s sqlStorage) FindAll() ([]*Domain, error) {
	var domains []*Domain

	if err := s.db.Preload("User").Find(&domains).Error; err != nil {
		return nil, err
	}

	return domains, nil
}

func (s sqlStorage) FindByUserID(userID int) ([]*Domain, error) {
	var domains []*Domain

	err := s.db.Preload("User").Where("user_id = ?", userID).Find(&domains).Error

	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (s sqlStorage) FindByID(id int) (*Domain, error) {
	var domain Domain

	if err := s.db.Preload("User").First(&domain, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &domain, nil
}

func (s sqlStorage) Save(domain *Domain) error {
	if domain.ID == 0 {
		return s.db.Create(domain).Error
	}

	return s.db.Save(domain).Error
}

func (s sqlStorage) Remove(domain *Domain) error {
	return s.db.Delete(domain).Error
}

func (s sqlStorage) UpdateCertData(domain *Domain, cert *agentintegration.Certificate) error {
	data, err := domain.GetData()

	if err != nil {
		return err
	}

	data.Cert = cert
	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	domain.Data = string(jsonData)

	return s.Save(domain)
}

func NewDomainStorage(db *gorm.DB) DomainStorage {
	return sqlStorage{
		db: db,
	}
}

func (Domain) TableName() string {
	return "certificate_monitor_sites"
}
