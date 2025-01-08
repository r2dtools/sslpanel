package storage

import "time"

type Server struct {
	ID           uint      `gorm:"AUTO_INCREMENT" gorm:"primary_key" json:"id"`
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
	FindByID(id int) (*Server, error)
	FindByGuid(guid string) (*Server, error)
	FindCountByIP(ipv4, ipv6 string, excludeIds []int) (int, error)
	Save(*Server) error
	Remove(*Server) error
}
