package storage

import (
	"backend/internal/app/panel/user/storage"
	"encoding/json"
	"errors"
	"time"

	"github.com/r2dtools/agentintegration"
)

type Domain struct {
	ID        uint         `gorm:"AUTO_INCREMENT" gorm:"primary_key" json:"id"`
	URL       string       `gorm:"size:64" json:"url"`
	Data      string       `gorm:"size:4096" json:"data"`
	UserID    uint         `json:"user_id"`
	User      storage.User `json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type DomainData struct {
	Cert *agentintegration.Certificate
}

func (d Domain) GetData() (DomainData, error) {
	var data DomainData

	if d.Data != "" {
		if err := json.Unmarshal([]byte(d.Data), &data); err != nil {
			return data, err
		}

		return data, nil
	}

	return data, errors.New("empty domain data")
}

type DomainStorage interface {
	FindAll() ([]*Domain, error)
	FindByUserID(userID int) ([]*Domain, error)
	FindByID(id int) (*Domain, error)
	Save(domain *Domain) error
	Remove(domain *Domain) error
	UpdateCertData(domain *Domain, cert *agentintegration.Certificate) error
}
