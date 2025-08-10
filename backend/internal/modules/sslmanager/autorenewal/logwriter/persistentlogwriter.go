package logwriter

import (
	"backend/internal/modules/sslmanager/autorenewal/logstorage"
)

type PersistentLogWriter struct {
	storage *logstorage.SqlRenewalLogStorage
}

func (w *PersistentLogWriter) WriteLog(serverID uint, successDomains []string, failedDomains map[string]error) error {
	records := []logstorage.RenewalLog{}

	for _, domainName := range successDomains {
		records = append(records, logstorage.RenewalLog{
			ServerID:   serverID,
			DomainName: domainName,
		})
	}

	for domainName, err := range failedDomains {
		records = append(records, logstorage.RenewalLog{
			ServerID:   serverID,
			DomainName: domainName,
			Error:      err.Error(),
		})
	}

	if len(records) == 0 {
		return nil
	}

	return w.storage.CreateLogs(records)
}

func CreatePersistentLogWriter(storage *logstorage.SqlRenewalLogStorage) *PersistentLogWriter {
	return &PersistentLogWriter{storage: storage}
}
