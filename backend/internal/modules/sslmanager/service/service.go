package service

import (
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/modules/sslmanager/agent"
	"errors"
	"fmt"

	serverAgent "backend/internal/pkg/agent"
	"backend/internal/pkg/certificate"

	"github.com/r2dtools/agentintegration"
)

var ErrServerNotFound = errors.New("server not found")

type CertificateService struct {
	serverStorage serverStorage.ServerStorage
}

func (s CertificateService) IssueCertificate(
	serverID int,
	certData agentintegration.CertificateIssueRequestData,
) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(serverID)

	if err != nil {
		return nil, err
	}

	return cAgent.Issue(certData)
}

func (s CertificateService) AssignCertificate(
	serverID int,
	requestData agentintegration.CertificateAssignRequestData,
) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(serverID)

	if err != nil {
		return nil, err
	}

	return cAgent.AssignCertificateToDomain(&requestData)
}

func (s CertificateService) UploadCertificate(
	serverID int,
	requestData agentintegration.CertificateUploadRequestData,
) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(serverID)

	if err != nil {
		return nil, err
	}

	return cAgent.Upload(&requestData)
}

func (s CertificateService) UploadCertificateToStorage(
	serverID int,
	requestData agentintegration.CertificateUploadRequestData,
) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(serverID)

	if err != nil {
		return nil, err
	}

	return cAgent.UploadPemCertificateToStorage(&requestData)
}

func (s CertificateService) DownloadCertificateFromStorage(serverID int, certName string) (*agentintegration.CertificateDownloadResponseData, error) {
	cAgent, err := s.getCertificateAgent(serverID)

	if err != nil {
		return nil, err
	}

	return cAgent.DownloadtStorageCertificate(certName)
}

func (s CertificateService) GetStorageCertNameList(serverID int) ([]string, error) {
	cAgent, err := s.getCertificateAgent(serverID)

	if err != nil {
		return nil, err
	}

	return cAgent.GetStorageCertNameList()
}

func (s CertificateService) GetStorageCertificate(serverID int, certName string) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(serverID)

	if err != nil {
		return nil, err
	}

	return cAgent.GetStorageCertificate(certName)
}

func (s CertificateService) RemoveCertificateFromStorage(serverID int, certName string) error {
	cAgent, err := s.getCertificateAgent(serverID)

	if err != nil {
		return err
	}

	return cAgent.RemoveCertificatefromStorage(certName)
}

func (s CertificateService) CreateSelfSignCertificate(serverID int, requestData SelfSignedCertificateRequest) (*agentintegration.Certificate, error) {
	cAgent, err := s.getCertificateAgent(serverID)

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

func (s CertificateService) getCertificateAgent(serverID int) (*agent.CertificateAgent, error) {
	server, err := s.serverStorage.FindByID(serverID)

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
	)

	if err != nil {
		return nil, err
	}

	return agent.NewCertificateAgent(sAgent), nil
}

func NewCertificateService(serverStorage serverStorage.ServerStorage) CertificateService {
	return CertificateService{
		serverStorage: serverStorage,
	}
}
