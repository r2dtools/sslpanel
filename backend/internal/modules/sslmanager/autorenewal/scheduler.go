package autorenewal

import (
	"backend/config"
	"backend/internal/pkg/logger"
	"fmt"
	"time"
)

type Scheduler struct {
	config  *config.Config
	logger  logger.Logger
	manager AutoRenewalManager
}

func (s Scheduler) Run() {
	limiter := make(chan struct{}, 1)
	tick := time.Tick(s.config.CertRenewalInterval)

	for t := range tick {
		select {
		case limiter <- struct{}{}:
			s.logger.Debug(fmt.Sprintf("start renewal: %v", t))
			go s.manager.Run(limiter)
		default:
			s.logger.Warning(fmt.Sprintf("renewal is in progress: %v", t))
		}
	}
}

func CreateScheduler(config *config.Config, logger logger.Logger, manager AutoRenewalManager) Scheduler {
	return Scheduler{
		config:  config,
		logger:  logger,
		manager: manager,
	}
}
