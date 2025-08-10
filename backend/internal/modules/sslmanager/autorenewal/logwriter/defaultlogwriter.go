package logwriter

import (
	"backend/internal/pkg/logger"
	"fmt"
)

type DefaultLogWriter struct {
	logger logger.Logger
}

func (s DefaultLogWriter) WriteLog(serverID uint, successDomains []string, failedDomains map[string]error) error {
	for _, domainName := range successDomains {
		s.logger.Info(fmt.Sprintf("certificate successfully renewed, server: %d, domain: %s", serverID, domainName))
	}

	for domainName, err := range failedDomains {
		s.logger.Error(fmt.Sprintf("certificate renewal failed, server: %d, domain: %s, error: %v", serverID, domainName, err))
	}

	return nil
}

func CreateDefaultLogWriter(logger logger.Logger) DefaultLogWriter {
	return DefaultLogWriter{logger: logger}
}
