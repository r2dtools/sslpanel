package storage

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const guidPrefix = "server_guid_prefix"

type Server struct {
	ID           uint      `gorm:"AUTO_INCREMENT;primary_key" json:"id"`
	Guid         string    `gorm:"-" json:"guid"`
	Name         string    `gorm:"size:64" json:"name"`
	OsCode       string    `gorm:"size:64" json:"os_code"`
	OsVersion    string    `gorm:"size:64" json:"os_version"`
	Ipv4Address  string    `gorm:"size:64" json:"ipv4_address"`
	Ipv6Address  string    `gorm:"size:256" json:"ipv6_address"`
	AgentVersion string    `gorm:"size:64" json:"agent_version"`
	AgentPort    int       `json:"agent_port"`
	Token        string    `gorm:"size:512" json:"token"`
	IsActive     uint8     `json:"is_active"`
	IsRegistered uint8     `json:"is_registered"`
	AccountID    uint      `json:"account_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ServerStorage interface {
	FindAllByAccountID(accountID int) ([]Server, error)
	FindAll() ([]Server, error)
	FindByID(id int) (*Server, error)
	FindByGuid(guid string) (*Server, error)
	FindCountByIP(ipv4, ipv6 string, excludeIds []int) (int, error)
	Save(*Server) error
	Remove(*Server) error
}

func GetServerIDByGUID(guid string) (int, error) {
	decodedGuid, err := base64.RawStdEncoding.DecodeString(guid)

	if err != nil {
		return 0, fmt.Errorf("invalid server guid: %v", err)
	}

	idStr := strings.TrimPrefix(string(decodedGuid), fmt.Sprintf("%s_", guidPrefix))
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return 0, fmt.Errorf("invalid server guid: %v", err)
	}

	return id, nil
}

func GetServerGUIDByID(id int) string {
	guidRaw := fmt.Sprintf("%s_%d", guidPrefix, id)

	return base64.RawStdEncoding.EncodeToString([]byte(guidRaw))
}
