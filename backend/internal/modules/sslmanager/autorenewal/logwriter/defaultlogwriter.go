package logwriter

import (
	"backend/internal/pkg/logger"
	"fmt"
)

type DefaultLogWriter struct {
	logger logger.Logger
}

func (s DefaultLogWriter) WriteLog(serverName string, successDomains []string, failedDomains map[string]error) error {
	for _, domainName := range successDomains {
		s.logger.Info(fmt.Sprintf("successfully renewed certificates, server: %s, domain: %s", serverName, domainName))
	}

	for domainName, err := range failedDomains {
		s.logger.Error(fmt.Sprintf("renewal failed, server: %s, domain: %s, error: %v", serverName, domainName, err))
	}

	return nil
}

func CreateDefaultLogWriter(logger logger.Logger) DefaultLogWriter {
	return DefaultLogWriter{logger: logger}
}
