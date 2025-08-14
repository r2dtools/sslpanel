package logstorage

import (
	serverStorage "backend/internal/app/panel/server/storage"
	"time"
)

type RenewalLog struct {
	ID         int `gorm:"AUTO_INCREMENT;primary_key"`
	ServerID   uint
	Server     serverStorage.Server
	DomainName string
	Error      string
	CreatedAt  time.Time
}

type RenewalLogStorage interface {
	CreateLogs(logs []RenewalLog) error
	FindLatestLogs(serverGuid string) ([]RenewalLog, error)
}
