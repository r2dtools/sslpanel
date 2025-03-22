package service

import (
	"backend/internal/app/panel/domain/service"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/modules/sslmanager/agent"
	"errors"
	"fmt"

	serverAgent "backend/internal/pkg/agent"
	"backend/internal/pkg/certificate"
	"backend/internal/pkg/logger"

	"github.com/r2dtools/agentintegration"
)

type ErrAgentCommon struct {
	message string
}

func (err ErrAgentCommon) Error() string {
	return err.message
}

var ErrServerNotFound = errors.New("server not found")

type CertificateService struct {
	serverStorage serverStorage.ServerStorage
	logger        logger.Logger
}

func (s CertificateService) IssueCertificate(request CertificateIssueRequest) (*service.DomainCertificate, error) {
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return nil, err
	}

	cert, err := cAgent.Issue(agentintegration.CertificateIssueRequestData{
		Email:            request.Email,
		ServerName:       request.DomainName,
		WebServer:        request.WebServer,
		ChallengeType:    request.ChallengeType,
		Subjects:         request.Subjects,
		AdditionalParams: request.AdditionalParams,
		Assign:           request.Assign,
	})

	if err != nil {
		return nil, ErrAgentCommon{message: err.Error()}
	}

	return service.CreateCertificate(cert), nil
}

func (s CertificateService) AssignCertificate(request AssignCertificateRequest) (*service.DomainCertificate, error) {
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return nil, err
	}

	cert, err := cAgent.AssignCertificateToDomain(&agentintegration.CertificateAssignRequestData{
		ServerName: request.DomainName,
		WebServer:  request.WebServer,
		CertName:   request.CertName,
	})

	if err != nil {
		return nil, err
	}

	return service.CreateCertificate(cert), nil
}

func (s CertificateService) UploadCertificate(
	guid string,
	requestData agentintegration.CertificateUploadRequestData,
) (*service.DomainCertificate, error) {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return nil, err
	}

	cert, err := cAgent.Upload(&requestData)

	if err != nil {
		return nil, err
	}

	return service.CreateCertificate(cert), nil
}

func (s CertificateService) UploadCertificateToStorage(request CertificateUploadToStorageRequest) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return nil, err
	}

	requestData := agentintegration.CertificateUploadRequestData{
		CertName:       request.CertName,
		PemCertificate: request.PemCertificate,
	}

	return cAgent.UploadPemCertificateToStorage(&requestData)
}

func (s CertificateService) DownloadCertificateFromStorage(guid string, certName string) (*agentintegration.CertificateDownloadResponseData, error) {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return nil, err
	}

	return cAgent.DownloadtStorageCertificate(certName)
}

func (s CertificateService) GetStorageCertificates(request CertificatesRequest) (map[string]Certificate, error) {
	result := map[string]Certificate{}
	cAgent, err := s.getCertificateAgent(request.Guid)

	if err != nil {
		return result, err
	}

	certsMap, err := cAgent.GetStorageCertificates()

	if err != nil {
		return result, ErrAgentCommon{message: err.Error()}
	}

	for name, cert := range certsMap {
		result[name] = createCertificate(cert)
	}

	return result, nil
}

func (s CertificateService) GetStorageCertificate(guid string, certName string) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return nil, err
	}

	return cAgent.GetStorageCertificate(certName)
}

func (s CertificateService) RemoveCertificateFromStorage(guid string, certName string) error {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return err
	}

	return cAgent.RemoveCertificatefromStorage(certName)
}

func (s CertificateService) GetCommonDirStatus(request CommonDirStatusRequest) (CommonDirStatusResponse, error) {
	var response CommonDirStatusResponse
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return response, err
	}

	agentResponse, err := cAgent.GetCommonDirStatus(agentintegration.CommonDirStatusRequestData{
		WebServer:  request.WebServer,
		ServerName: request.DomainName,
	})

	if err != nil {
		return response, err
	}

	response.Status = agentResponse.Status

	return response, nil
}

func (s CertificateService) ChangeCommonDirStatus(request ChangeCommonDirStatusRequest) error {
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return err
	}

	err = cAgent.ChangeCommonDirStatus(agentintegration.CommonDirChangeStatusRequestData{
		WebServer:  request.WebServer,
		ServerName: request.DomainName,
		Status:     request.Status,
	})

	if err != nil {
		return ErrAgentCommon{
			message: err.Error(),
		}
	}

	return nil
}

func (s CertificateService) CreateSelfSignCertificate(request SelfSignedCertificateRequest) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return nil, err
	}

	certData := certificate.SelfSignCertificateData{
		CertName:     request.CertName,
		CommonName:   request.CommonName,
		Email:        request.Email,
		Country:      request.Country,
		Province:     request.Province,
		Locality:     request.Locality,
		Organization: request.Organization,
		AltNames:     request.AltNames,
	}
	certPem, err := certificate.CreateSelfSignedCertificate(certData)

	if err != nil {
		return nil, fmt.Errorf("could not generate self-signed certificate: %v", err)
	}

	requestData := agentintegration.CertificateUploadRequestData{
		CertName:       request.CertName,
		PemCertificate: certPem,
	}

	return cAgent.UploadPemCertificateToStorage(&requestData)
}

func (s CertificateService) getCertificateAgent(guid string) (*agent.CertificateAgent, error) {
	server, err := s.serverStorage.FindByGuid(guid)

	if err != nil {
		return nil, err
	}

	if server == nil {
		return nil, ErrServerNotFound
	}

	sAgent, err := serverAgent.NewAgent(
		server.Ipv4Address,
		server.Ipv6Address,
		server.Token,
		server.AgentPort,
		s.logger,
	)

	if err != nil {
		return nil, err
	}

	return agent.NewCertificateAgent(sAgent), nil
}

func NewCertificateService(serverStorage serverStorage.ServerStorage, logger logger.Logger) CertificateService {
	return CertificateService{
		serverStorage: serverStorage,
		logger:        logger,
	}
}

func createCertificate(cert *agentintegration.Certificate) Certificate {
	return Certificate{
		CN:           cert.CN,
		ValidFrom:    cert.ValidFrom,
		ValidTo:      cert.ValidTo,
		DNSNames:     cert.DNSNames,
		Emails:       cert.EmailAddresses,
		Organization: cert.Organization,
		Province:     cert.Province,
		Country:      cert.Country,
		Locality:     cert.Locality,
		IsCA:         cert.IsCA,
		IsValid:      cert.IsValid,
		Issuer: Issuer{
			CN:           cert.Issuer.CN,
			Organization: cert.Issuer.Organization,
		},
	}
}
