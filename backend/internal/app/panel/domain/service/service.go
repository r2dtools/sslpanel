package service

import (
	"backend/config"
	"backend/internal/app/panel/domain/dto"
	"backend/internal/app/panel/domain/provider"
	"backend/internal/app/panel/domain/storage"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/pkg/agent"
	"backend/internal/pkg/logger"
	"errors"
	"fmt"

	"github.com/r2dtools/agentintegration"
)

var ErrServerNotFound = errors.New("server not found")
var ErrDomainNotFound = errors.New("domain not found")
var ErrAgentConnection = errors.New("failed to connect to the server agent")

type DomainService struct {
	config         *config.Config
	settingStorage storage.DomainSettingStorage
	serverStorage  serverStorage.ServerStorage
	domainProvider provider.DomainProvider
	logger         logger.Logger
}

func (s DomainService) GetDomain(request DomainRequest) (dto.Domain, error) {
	var rDomain dto.Domain

	serverModel, err := s.serverStorage.FindByGuid(request.ServerGuid)

	if err != nil {
		return rDomain, err
	}

	if serverModel == nil || serverModel.AccountID != uint(request.AccountID) {
		return rDomain, ErrServerNotFound
	}

	domains, err := s.domainProvider.GetServerDomains(request.ServerGuid)

	if err == provider.ErrServerNotFound {
		return rDomain, ErrServerNotFound
	}

	for _, domain := range domains {
		if domain.ServerName == request.DomainName {
			if request.WebServer == "" || request.WebServer == domain.WebServer {
				return domain, nil
			}
		}
	}

	return rDomain, ErrDomainNotFound
}

func (s DomainService) GetDomainConfig(request DomainConfigRequest) (string, error) {
	serverModel, err := s.serverStorage.FindByGuid(request.ServerGuid)

	if err != nil {
		return "", err
	}

	if serverModel == nil || serverModel.AccountID != uint(request.AccountID) {
		return "", ErrServerNotFound
	}

	nAgent, err := s.getServerAgent(serverModel)

	if err != nil {
		return "", err
	}

	config, err := nAgent.GetVhostConfig(agentintegration.VirtualHostConfigRequestData{
		WebServer:  request.WebServer,
		ServerName: request.DomainName,
	})

	if err != nil {
		return "", err
	}

	return config.Content, nil
}

func (s DomainService) FindDomainSettings(request DomainSettingsRequest) ([]DomainSetting, error) {
	settings := []DomainSetting{}

	serverModel, err := s.serverStorage.FindByGuid(request.ServerGuid)

	if err != nil {
		return settings, err
	}

	if serverModel == nil || serverModel.AccountID != uint(request.AccountID) {
		return settings, ErrServerNotFound
	}

	settingModels, err := s.settingStorage.FindAllByDomain(request.DomainName, request.ServerGuid)

	if err != nil {
		return settings, fmt.Errorf("error while searching settings for domain %s", request.DomainName)
	}

	for _, settingModel := range settingModels {
		settings = append(settings, createDomainSetting(settingModel))
	}

	return settings, nil
}

func (s DomainService) ChangeDomainSettings(request ChangeDomainSettingRequest) error {
	serverModel, err := s.serverStorage.FindByGuid(request.ServerGuid)

	if err != nil {
		return err
	}

	if serverModel == nil || serverModel.AccountID != uint(request.AccountID) {
		return ErrServerNotFound
	}

	settingModel, err := s.settingStorage.FindByDomain(request.DomainName, request.ServerGuid, request.SettingName)

	if err != nil {
		return fmt.Errorf("error while searching setting %s for domain %s", request.SettingName, request.DomainName)
	}

	if settingModel == nil {
		err = s.settingStorage.Create(
			request.DomainName,
			request.ServerGuid,
			request.SettingName,
			request.SettingValue,
		)
	} else {
		settingModel.SettingValue = request.SettingValue
		err = s.settingStorage.Save(settingModel)
	}

	if err != nil {
		return fmt.Errorf("failed to change setting %s: %v", request.SettingName, err)
	}

	return nil
}

func (s DomainService) getServerAgent(server *serverStorage.Server) (*agent.Agent, error) {
	return agent.NewAgent(
		server.Ipv4Address,
		server.Ipv6Address,
		server.Token,
		server.AgentPort,
		s.logger,
	)
}

func createDomainSetting(settingModel storage.DomainSetting) DomainSetting {
	return DomainSetting{
		ID:           settingModel.ID,
		SettingName:  settingModel.SettingName,
		SettingValue: settingModel.SettingValue,
	}
}

func NewDomainService(
	config *config.Config,
	settingStorage storage.DomainSettingStorage,
	serverStorage serverStorage.ServerStorage,
	domainProvider provider.DomainProvider,
	logger logger.Logger,
) DomainService {
	return DomainService{
		config:         config,
		settingStorage: settingStorage,
		serverStorage:  serverStorage,
		domainProvider: domainProvider,
		logger:         logger,
	}
}
