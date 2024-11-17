package service

import (
	"backend/config"
	"backend/internal/app/panel/server/storage"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/pkg/agent"
	"backend/internal/pkg/logger"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"github.com/r2dtools/agentintegration"
)

const (
	serverTokenLength = 24
)

var ErrServerNotFound = errors.New("server not found")

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

func (s ServerService) FindServerByID(id int) (*Server, error) {
	serverModel, err := s.serverStorage.FindByID(id)

	if err != nil {
		return nil, fmt.Errorf("could not find server with ID %d: %w", id, err)
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
		return errors.New("ipv4 or ipv6 address must be specified")
	}

	count, err := s.serverStorage.FindCountByIP(request.Ipv4Address, request.Ipv6Address, nil)

	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("server with the specified ipv4 or ipv6 address already exists")
	}

	serverModel := &serverStorage.Server{
		Name:        request.Name,
		Ipv4Address: request.Ipv4Address,
		Ipv6Address: request.Ipv6Address,
		AgentPort:   request.AgentPort,
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
		return errors.New("ipv4 or ipv6 address must be specified")
	}

	count, err := s.serverStorage.FindCountByIP(request.Ipv4Address, request.Ipv6Address, []int{request.ID})

	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("server with the specified ipv4 or ipv6 address already exists")
	}

	serverModel, err := s.serverStorage.FindByID(request.ID)

	if err != nil {
		return err
	}

	if serverModel == nil {
		return ErrServerNotFound
	}

	serverModel.Name = request.Name
	serverModel.OsCode = request.OsCode
	serverModel.OsVersion = request.OsVersion
	serverModel.Ipv4Address = request.Ipv4Address
	serverModel.Ipv6Address = request.Ipv6Address
	serverModel.IsActive = uint8(request.IsActive)
	serverModel.IsRegistered = uint8(request.IsRegistered)

	agentPort := request.AgentPort

	if agentPort == 0 {
		agentPort = s.config.AgentPort
	}

	serverModel.AgentPort = agentPort

	return s.serverStorage.Save(serverModel)
}

func (s ServerService) RefreshServer(serverID int) (*Server, *agentintegration.ServerData, error) {
	serverModel, err := s.serverStorage.FindByID(serverID)

	if err != nil {
		return nil, nil, err
	}

	if serverModel == nil {
		return nil, nil, ErrServerNotFound
	}

	nAgent, err := getServerAgent(serverModel)

	if err != nil {
		return nil, nil, err
	}

	refData, err := nAgent.Refresh()
	connErr := agent.ConnectionError{}

	if err != nil {
		// If connection failed with server agent then change agent status to "inactive"
		if errors.As(err, &connErr) {
			serverModel.IsActive = 0
			err = s.serverStorage.Save(serverModel)

			s.logger.Error(err.Error())
		}

		return nil, nil, err
	}

	serverModel.OsCode = refData.Platform
	serverModel.OsVersion = refData.PlatformVersion
	serverModel.AgentVersion = refData.AgentVersion
	serverModel.IsRegistered = 1
	serverModel.IsActive = 1

	err = s.serverStorage.Save(serverModel)

	if err != nil {
		return nil, nil, err
	}

	return createServer(serverModel), refData, nil
}

func (s ServerService) FindServerVhosts(serverID int) ([]agentintegration.VirtualHost, error) {
	serverModel, err := s.serverStorage.FindByID(serverID)

	if err != nil {
		return nil, err
	}

	if serverModel == nil {
		return nil, ErrServerNotFound
	}

	nAgent, err := getServerAgent(serverModel)

	if err != nil {
		return nil, err
	}

	vhosts, err := nAgent.GetVhosts()

	if err != nil {
		return nil, err
	}

	return filterVhosts(vhosts), nil
}

func (s ServerService) GetVhostCertificate(serverID int, vhostName string) (*agentintegration.Certificate, error) {
	serverModel, err := s.serverStorage.FindByID(serverID)

	if err != nil {
		return nil, err
	}

	if serverModel == nil {
		return nil, ErrServerNotFound
	}

	nAgent, err := getServerAgent(serverModel)

	if err != nil {
		return nil, err
	}

	return nAgent.GetVhostCertificate(vhostName)
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

func getServerAgent(server *serverStorage.Server) (*agent.Agent, error) {
	return agent.NewAgent(
		server.Ipv4Address,
		server.Ipv6Address,
		server.Token,
		server.AgentPort,
	)
}

func filterVhosts(vhosts []agentintegration.VirtualHost) []agentintegration.VirtualHost {
	var rVhosts []agentintegration.VirtualHost

	for _, vhost := range vhosts {
		serverName := strings.Trim(vhost.ServerName, ".")
		serverNameParts := strings.Split(serverName, ".")

		// skip vhost names like 'domain'
		if len(serverNameParts) > 1 {
			rVhosts = append(rVhosts, vhost)
		}
	}

	return rVhosts
}

func NewServerService(config *config.Config, serverStorage storage.ServerStorage, logger logger.Logger) ServerService {
	return ServerService{
		config:        config,
		serverStorage: serverStorage,
		logger:        logger,
	}
}
