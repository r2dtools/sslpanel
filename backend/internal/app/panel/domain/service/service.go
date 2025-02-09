package service

import (
	"backend/config"
	"backend/internal/app/panel/domain/storage"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/pkg/agent"
	"backend/internal/pkg/logger"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/r2dtools/agentintegration"
)

var ErrServerNotFound = errors.New("server not found")
var ErrDomainNotFound = errors.New("domain not found")
var ErrAgentConnection = errors.New("failed to connect to the server agent")

type ErrAgentCommon struct {
	message string
}

func (err ErrAgentCommon) Error() string {
	return err.message
}

type DomainService struct {
	config         *config.Config
	settingStorage storage.DomainSettingStorage
	serverStorage  serverStorage.ServerStorage
	logger         logger.Logger
}

func (s DomainService) GetDomain(request DomainRequest) (*Domain, error) {
	serverModel, err := s.serverStorage.FindByGuid(request.ServerGuid)

	if err != nil {
		return nil, err
	}

	if serverModel == nil {
		return nil, ErrServerNotFound
	}

	nAgent, err := s.getServerAgent(serverModel)

	if err != nil {
		return nil, err
	}

	vhosts, err := nAgent.GetVhosts()

	if err != nil {
		return nil, err
	}

	var dVhost *agentintegration.VirtualHost

	for _, vhost := range vhosts {
		if vhost.ServerName == request.DomainName {
			if request.WebServer == "" || request.WebServer == vhost.WebServer {
				dVhost = &vhost

				break
			}
		}
	}

	if dVhost == nil {
		return nil, ErrDomainNotFound
	}

	domain := CreateDomain(dVhost)

	if domain == nil {
		return nil, ErrDomainNotFound
	}

	return domain, nil
}

func (s DomainService) GetDomainConfig(request DomainConfigRequest) (string, error) {
	serverModel, err := s.serverStorage.FindByGuid(request.ServerGuid)

	if err != nil {
		return "", err
	}

	if serverModel == nil {
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
		return "", ErrAgentCommon{message: err.Error()}
	}

	return config.Content, nil
}

func (s DomainService) FindDomainSettings(domainName string, serverGuid string) ([]DomainSetting, error) {
	var settings []DomainSetting
	settingModels, err := s.settingStorage.FindByDomain(domainName, serverGuid)

	if err != nil {
		return settings, fmt.Errorf("could not find domain settings for domain %s", domainName)
	}

	for _, settingModel := range settingModels {
		settings = append(settings, createDomainSetting(settingModel))
	}

	return settings, nil
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

func CreateDomain(vhost *agentintegration.VirtualHost) *Domain {
	serverName := strings.Trim(vhost.ServerName, ".")
	serverNameParts := strings.Split(serverName, ".")

	// skip vhost names like 'domain'
	if len(serverNameParts) <= 1 {
		return nil
	}

	var addresses []DomainAddress

	for _, address := range vhost.Addresses {
		port, err := strconv.Atoi(address.Port)

		if err != nil {
			continue
		}

		addresses = append(addresses, DomainAddress{
			IsIpv6: address.IsIpv6,
			Host:   address.Host,
			Port:   port,
		})
	}

	return &Domain{
		FilePath:    vhost.FilePath,
		ServerName:  vhost.ServerName,
		DocRoot:     vhost.DocRoot,
		WebServer:   vhost.WebServer,
		Aliases:     vhost.Aliases,
		Ssl:         vhost.Ssl,
		Addresses:   addresses,
		Certificate: CreateCertificate(vhost.Certificate),
	}
}

func createDomainSetting(settingModel storage.DomainSetting) DomainSetting {
	return DomainSetting{
		ID:           settingModel.ID,
		SettingName:  settingModel.SettingName,
		SettingValue: settingModel.SettingValue,
	}
}

func CreateCertificate(cert *agentintegration.Certificate) *DomainCertificate {
	if cert == nil {
		return nil
	}

	return &DomainCertificate{
		CN:             cert.CN,
		ValidFrom:      cert.ValidFrom,
		ValidTo:        cert.ValidTo,
		DNSNames:       cert.DNSNames,
		EmailAddresses: cert.EmailAddresses,
		Organization:   cert.Organization,
		Country:        cert.Country,
		Locality:       cert.Locality,
		Province:       cert.Province,
		IsValid:        cert.IsValid,
		IsCA:           cert.IsCA,
		Issuer:         Issuer(cert.Issuer),
	}
}

func NewDomainService(
	config *config.Config,
	settingStorage storage.DomainSettingStorage,
	serverStorage serverStorage.ServerStorage,
	logger logger.Logger,
) DomainService {
	return DomainService{
		config:         config,
		settingStorage: settingStorage,
		serverStorage:  serverStorage,
		logger:         logger,
	}
}
