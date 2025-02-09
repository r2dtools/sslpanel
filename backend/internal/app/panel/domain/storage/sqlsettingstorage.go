package storage

import (
	serverStorage "backend/internal/app/panel/server/storage"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type sqlSettingStorage struct {
	db *gorm.DB
}

func (s sqlSettingStorage) FindByID(id int) (*DomainSetting, error) {
	var setting DomainSetting
	err := s.db.First(&setting, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not find domain setting with ID %d: %v", id, err)
	}

	return &setting, nil
}

func (s sqlSettingStorage) FindByDomain(domainName string, serverGuid string) ([]DomainSetting, error) {
	var settings []DomainSetting
	serverId, err := serverStorage.GetServerIDByGUID(serverGuid)

	if err != nil {
		return settings, err
	}

	err = s.db.Where("server_id = ?", serverId).Where("domain_name = ?", domainName).Find(&settings).Error

	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (s sqlSettingStorage) Save(setting *DomainSetting) error {
	if setting.ID == 0 {
		return s.db.Create(setting).Error
	}

	return s.db.Save(setting).Error
}

func (s sqlSettingStorage) Remove(setting *DomainSetting) error {
	err := s.db.Delete(setting).Error

	if err != nil {
		return fmt.Errorf("failed to delete domain setting with ID %d: %v", setting.ID, err)
	}

	return nil
}

func NewDomainSettingSqlStorage(db *gorm.DB) DomainSettingStorage {
	return sqlSettingStorage{db: db}
}

func (*DomainSetting) TableName() string {
	return "domain_settings"
}
