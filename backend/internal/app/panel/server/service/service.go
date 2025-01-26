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
	"strconv"
	"strings"

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

var ErrServerNotFound = errors.New("server not found")
var ErrDomainNotFound = errors.New("domain not found")
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
			s.serverStorage.Save(serverModel)

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

func (s ServerService) GetServerDomain(guid string, domainName string) (*Domain, error) {
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

	vhosts, err := nAgent.GetVhosts()

	if err != nil {
		return nil, err
	}

	var dVhost *agentintegration.VirtualHost

	for _, vhost := range vhosts {
		if vhost.ServerName == domainName {
			dVhost = &vhost

			break
		}
	}

	if dVhost == nil {
		return nil, ErrDomainNotFound
	}

	domain := createDomain(dVhost)

	if domain == nil {
		return nil, ErrDomainNotFound
	}

	return domain, nil
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

func (s ServerService) GetVhostCertificate(serverID int, vhostName string) (*agentintegration.Certificate, error) {
	serverModel, err := s.serverStorage.FindByID(serverID)

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

func createDomains(vhosts []agentintegration.VirtualHost) []Domain {
	var domains []Domain

	for _, vhost := range vhosts {
		domain := createDomain(&vhost)

		if domain == nil {
			continue
		}

		domains = append(domains, *domain)
	}

	return domains
}

func createDomain(vhost *agentintegration.VirtualHost) *Domain {
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

func NewServerService(config *config.Config, serverStorage storage.ServerStorage, logger logger.Logger) ServerService {
	return ServerService{
		config:        config,
		serverStorage: serverStorage,
		logger:        logger,
	}
}
