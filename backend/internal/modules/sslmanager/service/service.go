package service

import (
	"backend/internal/app/panel/server/service"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/modules/sslmanager/agent"
	"errors"
	"fmt"

	serverAgent "backend/internal/pkg/agent"
	"backend/internal/pkg/certificate"
	"backend/internal/pkg/logger"

	"github.com/r2dtools/agentintegration"
)

type ErrAgentCertificate struct {
	message string
}

func (e ErrAgentCertificate) Error() string {
	return e.message
}

var ErrServerNotFound = errors.New("server not found")

type CertificateService struct {
	serverStorage serverStorage.ServerStorage
	logger        logger.Logger
}

func (s CertificateService) IssueCertificate(
	guid string,
	certIssueRequest CertificateIssueRequest,
) (*service.DomainCertificate, error) {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return nil, err
	}

	cert, err := cAgent.Issue(agentintegration.CertificateIssueRequestData(certIssueRequest))

	if err != nil {
		return nil, ErrAgentCertificate{message: err.Error()}
	}

	return service.CreateCertificate(cert), nil
}

func (s CertificateService) AssignCertificate(
	guid string,
	requestData agentintegration.CertificateAssignRequestData,
) (*service.DomainCertificate, error) {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return nil, err
	}

	cert, err := cAgent.AssignCertificateToDomain(&requestData)

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

func (s CertificateService) UploadCertificateToStorage(
	guid string,
	requestData agentintegration.CertificateUploadRequestData,
) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return nil, err
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

func (s CertificateService) GetStorageCertNameList(guid string) ([]string, error) {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return nil, err
	}

	return cAgent.GetStorageCertNameList()
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

func (s CertificateService) CreateSelfSignCertificate(guid string, requestData SelfSignedCertificateRequest) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return nil, err
	}

	certData := certificate.SelfSignCertificateData{
		CertName:     requestData.CertName,
		CommonName:   requestData.CommonName,
		Email:        requestData.Email,
		Country:      requestData.Country,
		Province:     requestData.Province,
		Locality:     requestData.Locality,
		Organization: requestData.Organization,
		AltNames:     requestData.AltNames,
	}
	certPem, err := certificate.CreateSelfSignedCertificate(certData)

	if err != nil {
		return nil, fmt.Errorf("could not generate self-signed certificate: %v", err)
	}

	var uploadCertRequestData agentintegration.CertificateUploadRequestData
	uploadCertRequestData.PemCertificate = certPem
	uploadCertRequestData.CertName = requestData.CertName

	return cAgent.UploadPemCertificateToStorage(&uploadCertRequestData)
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
