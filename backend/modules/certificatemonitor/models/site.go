package models

import (
	"backend/db"
	"backend/models"
	"encoding/json"

	"github.com/r2dtools/agentintegration"
)

// SiteData site data
type SiteData struct {
	Cert *agentintegration.Certificate
}

// Site model
type Site struct {
	models.Model

	URL    string      `gorm:"size:64" json:"url"`
	Data   string      `gorm:"size:4096" json:"data"`
	UserID uint        `json:"user_id"`
	User   models.User `json:"-"`
}

// GetUserSites returns all sites by user ID
func GetUserSites(userID uint) ([]*Site, error) {
	var sites []*Site

	err := db.GetDB().Where("user_id = ?", userID).Find(&sites).Error

	if err != nil {
		return nil, err
	}

	return sites, nil
}

// SaveSite creates new site record
func SaveSite(s *Site) error {
	var err error

	if db.GetDB().NewRecord(s) {
		err = db.GetDB().Create(s).Error
	} else {
		err = db.GetDB().Save(s).Error
	}

	if err != nil {
		return err
	}

	return nil
}

// Remove removes site
func (s *Site) Remove() error {
	err := db.GetDB().Delete(s).Error

	if err != nil {
		return err
	}

	return nil
}

// GetUser returns site user
func (s *Site) GetUser() (*models.User, error) {
	var user models.User
	assoc := db.GetDB().Model(s).Association("User").Find(&user)

	if assoc.Error != nil {
		return nil, assoc.Error
	}

	return &user, nil
}

// UpdateCertData updates site certificate data
func (s *Site) UpdateCertData(cert *agentintegration.Certificate) error {
	sData, err := s.GetData()

	if err != nil {
		return err
	}

	sData.Cert = cert
	jsonData, err := json.Marshal(sData)

	if err != nil {
		return err
	}

	s.Data = string(jsonData)

	if err = SaveSite(s); err != nil {
		return err
	}

	return nil
}

// GetData returns site data
func (s *Site) GetData() (*SiteData, error) {
	var sData SiteData

	if s.Data != "" {
		if err := json.Unmarshal([]byte(s.Data), &sData); err != nil {
			return nil, err
		}
	}

	return &sData, nil
}

// GetSiteByID returns site by ID
func GetSiteByID(id int) (*Site, error) {
	var site Site

	if err := db.GetDB().First(&site, id).Error; err != nil {
		return nil, err
	}

	return &site, nil
}

// GetAllSites loads all sites
func GetAllSites() ([]*Site, error) {
	var sites []*Site

	if err := db.GetDB().Find(&sites).Error; err != nil {
		return nil, err
	}

	return sites, nil
}

// TableName overwrite table name
func (Site) TableName() string {
	return "certificate_monitor_sites"
}
