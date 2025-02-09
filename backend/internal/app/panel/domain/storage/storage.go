package storage

import "time"

type DomainSetting struct {
	ID           int       `gorm:"AUTO_INCREMENT" gorm:"primary_key"`
	ServerId     int       `gorm:"-"`
	DomainName   string    `gorm:"size:64"`
	SettingName  string    `gorm:"size:64"`
	SettingValue string    `gorm:"size:64"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DomainSettingStorage interface {
	FindByID(id int) (*DomainSetting, error)
	FindByDomain(domainName string, serverId string) ([]DomainSetting, error)
	Save(*DomainSetting) error
	Remove(*DomainSetting) error
}
