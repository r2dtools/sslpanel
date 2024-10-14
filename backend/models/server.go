package models

import (
	"backend/config"
	"backend/db"
	"crypto/md5"
	"crypto/rand"
	"fmt"

	"github.com/jinzhu/gorm"
)

const guid_prefix = "server_guid_prefix"

// Server model
type Server struct {
	Model

	Guid         string  `gorm:"-" json:"guid"`
	Name         *string `gorm:"size:64" json:"name"`
	OsCode       string  `gorm:"size:64" json:"os_code"`
	OsVersion    string  `gorm:"size:64" json:"os_version"`
	Ipv4Address  string  `gorm:"size:64" json:"ipv4_address"`
	Ipv6Address  string  `gorm:"size:256" json:"ipv6_address"`
	AgentVersion string  `gorm:"size:64" json:"agent_version"`
	AgentPort    int     `json:"agent_port"`
	Token        string  `gorm:"size:512" json:"token"`
	IsActive     uint8   `json:"is_active"`
	IsRegistered uint8   `json:"is_registered"`
	AccountID    uint    `json:"account_id"`
}

// GetAccountServers returns all servers by account ID
func GetAccountServers(accountID uint) ([]Server, error) {
	var servers []Server
	err := db.GetDB().Where("account_id = ?", accountID).Order("id desc").Find(&servers).Error
	if err != nil {
		return nil, err
	}

	return servers, nil
}

// GetAccountServerByID returns server with the specified ID and account
func GetAccountServerByID(id int, accountID uint) (*Server, error) {
	var server Server
	err := db.GetDB().Where("id = ? AND account_id = ?", id, accountID).Find(&server).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return &server, nil
}

// Remove removes server
func (s *Server) Remove() error {
	err := db.GetDB().Delete(s).Error
	if err != nil {
		return err
	}

	return nil
}

// GetServerByID returns server by ID
func GetServerByID(id int) (*Server, error) {
	var server Server
	if err := db.GetDB().First(&server, id).Error; err != nil {
		return nil, err
	}

	return &server, nil
}

// SaveServer creates new server record
func SaveServer(s *Server) error {
	var err error

	if db.GetDB().NewRecord(s) {
		if s.Token == "" {
			token, err := generateToken(24)
			if err != nil {
				return fmt.Errorf("could not generate server token: %v", err)
			}

			s.Token = token
		}

		if s.AgentPort == 0 {
			s.AgentPort = config.GetConfig().AgentPort
		}

		err = db.GetDB().Create(s).Error
	} else {
		err = db.GetDB().Save(s).Error
	}

	if err != nil {
		return err
	}

	return nil
}

// BeforeCreate is a hook executed before creating new Server
func (s *Server) BeforeCreate(tx *gorm.DB) error {
	var ips []string

	if s.Ipv4Address != "" {
		ips = append(ips, s.Ipv4Address)
	}

	if s.Ipv6Address != "" {
		ips = append(ips, s.Ipv6Address)
	}

	if len(ips) == 0 {
		return nil
	}

	db := tx.Where("ipv4_address = ?", ips[0])

	if len(ips) == 2 {
		db = tx.Or("ipv6_address = ?", ips[1])
	}

	var count int
	err := db.Table("servers").Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("server with the specified ipv4 or ipv6 address already exists")
	}

	return nil
}

func (s *Server) AfterFind(tx *gorm.DB) (err error) {
	guidRaw := fmt.Sprintf("%s_%d", guid_prefix, s.ID)
	guidMd5 := md5.Sum([]byte(guidRaw))
	s.Guid = fmt.Sprintf("%x", guidMd5)
	return
}

func generateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
