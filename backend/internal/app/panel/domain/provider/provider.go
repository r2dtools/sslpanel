package provider

import (
	"backend/internal/app/panel/domain/dto"
	"backend/internal/app/panel/domain/factory"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/pkg/agent"
	"backend/internal/pkg/logger"
	"errors"
)

var ErrServerNotFound = errors.New("server not found")

type DomainProvider struct {
	serverStorage serverStorage.ServerStorage
	logger        logger.Logger
}

func (p DomainProvider) GetServerDomains(serverGuid string) ([]dto.Domain, error) {
	serverModel, err := p.serverStorage.FindByGuid(serverGuid)

	if err != nil {
		return nil, err
	}

	if serverModel == nil {
		return nil, ErrServerNotFound
	}

	nAgent, err := agent.NewAgent(
		serverModel.Ipv4Address,
		serverModel.Ipv6Address,
		serverModel.Token,
		serverModel.AgentPort,
		p.logger,
	)

	if err != nil {
		return nil, err
	}

	vhosts, err := nAgent.GetVhosts()

	if err != nil {
		return nil, err
	}

	domains := []dto.Domain{}

	for _, vhost := range vhosts {
		domain := factory.CreateDomain(vhost)

		if domain == nil {
			continue
		}

		domains = append(domains, *domain)
	}

	return domains, nil
}

func CreateDomainProvider(serverStorage serverStorage.ServerStorage, logger logger.Logger) DomainProvider {
	return DomainProvider{
		serverStorage: serverStorage,
		logger:        logger,
	}
}
