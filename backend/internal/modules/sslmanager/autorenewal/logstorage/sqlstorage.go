package logstorage

import (
	"backend/internal/app/panel/server/storage"
	"time"

	"gorm.io/gorm"
)

type SqlRenewalLogStorage struct {
	db *gorm.DB
}

func (s *SqlRenewalLogStorage) CreateLogs(logs []RenewalLog) error {
	return s.db.Create(&logs).Error
}

func (s *SqlRenewalLogStorage) FindLatestLogs(serverGuid string) ([]RenewalLog, error) {
	logs := []RenewalLog{}
	serverId, err := storage.GetServerIDByGUID(serverGuid)

	if err != nil {
		return logs, err
	}

	lastWeek := time.Now().AddDate(0, 0, -7)
	err = s.db.Preload("Server").
		Where("created_at > ?", lastWeek).
		Where("server_id = ?", serverId).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, err
}

func CreateSqlRenewalLogStorage(db *gorm.DB) *SqlRenewalLogStorage {
	return &SqlRenewalLogStorage{db: db}
}
