package logstorage

import "gorm.io/gorm"

type SqlRenewalLogStorage struct {
	db *gorm.DB
}

func (s *SqlRenewalLogStorage) CreateLogs(logs []RenewalLog) error {
	return s.db.Create(&logs).Error
}

func (s *SqlRenewalLogStorage) FindAllLogs(limit int) ([]RenewalLog, error) {
	logs := []RenewalLog{}
	err := s.db.Order("create_at DESC").Limit(limit).Find(&logs).Error

	return logs, err
}

func CreateSqlRenewalLogStorage(db *gorm.DB) *SqlRenewalLogStorage {
	return &SqlRenewalLogStorage{db: db}
}
