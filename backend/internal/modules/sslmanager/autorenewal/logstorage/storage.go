package logstorage

import "time"

type RenewalLog struct {
	ID         int `gorm:"AUTO_INCREMENT;primary_key"`
	ServerID   uint
	DomainName string `gorm:"size:64"`
	Error      string `gorm:"size:1024"`
	CreatedAt  time.Time
}

type RenewalLogStorage interface {
	CreateLogs(logs []RenewalLog) error
	FindAllLogs(limit int) ([]RenewalLog, error)
}
