package service

import (
	"backend/config"
	"backend/internal/app/panel/domain/dto"
	domainFactory "backend/internal/app/panel/domain/factory"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/pkg/agent"
	"backend/internal/pkg/logger"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/r2dtools/agentintegration"
)

const (
	serverTokenLength = 24
)

type ErrServerService struct {
	message string
}

func (e ErrServerService) Error() string {
	return e.message
}

type ErrAgentCommon struct {
	message string
}

func (err ErrAgentCommon) Error() string {
	return err.message
}

var ErrServerNotFound = errors.New("server not found")
var ErrAgentConnection = errors.New("failed to connect to the server agent")

type ServerService struct {
	config        *config.Config
	serverStorage serverStorage.ServerStorage
	logger        logger.Logger
}

func (s ServerService) FindAccountServers(accountID int) ([]Server, error) {
	var servers []Server
	serverModels, err := s.serverStorage.FindAllByAccountID(accountID)

	if err != nil {
		return servers, fmt.Errorf("could not get account %d servers", accountID)
	}

	for _, serverModel := range serverModels {
		servers = append(servers, *createServer(&serverModel))
	}

	return servers, nil
}

func (s ServerService) GetServerDetailsByGuid(guid string) (*ServerDetails, error) {
	serverModel, err := s.serverStorage.FindByGuid(guid)

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

	refData, err := nAgent.Refresh()
	connErr := agent.ConnectionError{}

	if err != nil {
		s.logger.Debug(err.Error())

		if errors.As(err, &connErr) {
			serverModel.IsActive = 0
			s.serverStorage.Save(serverModel) // nolint:errcheck

			return nil, ErrAgentConnection
		}

		return nil, err
	}

	serverModel.OsCode = refData.Platform
	serverModel.OsVersion = refData.PlatformVersion
	serverModel.AgentVersion = refData.AgentVersion
	serverModel.IsActive = 1

	err = s.serverStorage.Save(serverModel)

	if err != nil {
		return nil, err
	}

	vhosts, err := nAgent.GetVhosts()

	if err != nil {
		return nil, err
	}

	serverDetails := ServerDetails{
		Server:         *createServer(serverModel),
		PlatformFamily: refData.PlatformFamily,
		Os:             refData.Os,
		Virtualization: refData.Virtualization,
		HostName:       refData.HostName,
		KernelVersion:  refData.KernelVersion,
		KernelArch:     refData.KernelArch,
		Uptime:         refData.Uptime,
		BootTime:       refData.BootTime,
		Domains:        createDomains(vhosts),
	}

	return &serverDetails, nil
}

func (s ServerService) FindServerByGuid(guid string) (*Server, error) {
	serverModel, err := s.serverStorage.FindByGuid(guid)

	if err != nil {
		return nil, err
	}

	if serverModel == nil {
		return nil, ErrServerNotFound
	}

	return createServer(serverModel), nil
}

func (s ServerService) RemoveServer(id int) error {
	serverModel, err := s.serverStorage.FindByID(id)

	if err != nil {
		return err
	}

	if serverModel == nil {
		return ErrServerNotFound
	}

	return s.serverStorage.Remove(serverModel)
}

func (s ServerService) AddServer(request NewServerRequest) error {
	if request.Ipv4Address == "" && request.Ipv6Address == "" {
		return ErrServerService{"ipv4 or ipv6 address must be specified"}
	}

	count, err := s.serverStorage.FindCountByIP(request.Ipv4Address, request.Ipv6Address, nil)

	if err != nil {
		return err
	}

	if count > 0 {
		return ErrServerService{"server with the specified ipv4 or ipv6 address already exists"}
	}

	serverModel := &serverStorage.Server{
		Name:        request.Name,
		Ipv4Address: request.Ipv4Address,
		Ipv6Address: request.Ipv6Address,
		AgentPort:   request.AgentPort,
		Token:       request.Token,
	}
	guid, err := generateToken(serverTokenLength)

	if err != nil {
		return err
	}

	serverModel.Guid = guid
	serverModel.AccountID = uint(request.AccountID)

	agentPort := request.AgentPort

	if agentPort == 0 {
		agentPort = s.config.AgentPort
	}

	serverModel.AgentPort = agentPort

	return s.serverStorage.Save(serverModel)
}

func (s ServerService) UpdateServer(request UpdateServerRequest) error {
	if request.Ipv4Address == "" && request.Ipv6Address == "" {
		return ErrServerService{"ipv4 or ipv6 address must be specified"}
	}

	count, err := s.serverStorage.FindCountByIP(request.Ipv4Address, request.Ipv6Address, []int{request.ID})

	if err != nil {
		return err
	}

	if count > 0 {
		return ErrServerService{"server with the specified ipv4 or ipv6 address already exists"}
	}

	serverModel, err := s.serverStorage.FindByID(request.ID)

	if err != nil {
		return err
	}

	if serverModel == nil {
		return ErrServerNotFound
	}

	serverModel.Name = request.Name
	serverModel.Ipv4Address = request.Ipv4Address
	serverModel.Ipv6Address = request.Ipv6Address
	serverModel.Token = request.Token

	agentPort := request.AgentPort

	if agentPort == 0 {
		agentPort = s.config.AgentPort
	}

	serverModel.AgentPort = agentPort

	return s.serverStorage.Save(serverModel)
}

func createServer(server *serverStorage.Server) *Server {
	return &Server{
		ID:           int(server.ID),
		Guid:         server.Guid,
		Name:         server.Name,
		OsCode:       server.OsCode,
		OsVersion:    server.OsVersion,
		Ipv4Address:  server.Ipv4Address,
		Ipv6Address:  server.Ipv6Address,
		AgentVersion: server.AgentVersion,
		AgentPort:    server.AgentPort,
		IsActive:     int(server.IsActive),
		IsRegistered: int(server.IsRegistered),
		AccountID:    int(server.AccountID),
		CreatedAt:    server.CreatedAt,
		Token:        server.Token,
	}
}

func generateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s ServerService) getServerAgent(server *serverStorage.Server) (*agent.Agent, error) {
	return agent.NewAgent(
		server.Ipv4Address,
		server.Ipv6Address,
		server.Token,
		server.AgentPort,
		s.logger,
	)
}

func createDomains(vhosts []agentintegration.VirtualHost) []dto.Domain {
	var domains []dto.Domain

	for _, vhost := range vhosts {
		domain := domainFactory.CreateDomain(vhost)

		if domain == nil {
			continue
		}

		domains = append(domains, *domain)
	}

	return domains
}

func NewServerService(config *config.Config, serverStorage serverStorage.ServerStorage, logger logger.Logger) ServerService {
	return ServerService{
		config:        config,
		serverStorage: serverStorage,
		logger:        logger,
	}
}
