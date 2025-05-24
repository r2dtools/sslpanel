package storage

import "time"

type DomainSetting struct {
	ID           int `gorm:"AUTO_INCREMENT;primary_key"`
	ServerId     int
	DomainName   string `gorm:"size:64"`
	SettingName  string `gorm:"size:64"`
	SettingValue string `gorm:"size:64"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type DomainSettingStorage interface {
	FindByID(id int) (*DomainSetting, error)
	FindAllByDomain(domainName string, serverId string) ([]DomainSetting, error)
	FindByDomain(domainName string, serverId string, settingName string) (*DomainSetting, error)
	Create(domainName string, serverGuid string, settingName string, settingValue string) error
	Save(*DomainSetting) error
	Remove(*DomainSetting) error
}
