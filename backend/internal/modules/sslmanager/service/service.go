package service

import (
	"backend/internal/app/panel/domain/dto"
	domainFactory "backend/internal/app/panel/domain/factory"
	domainStorage "backend/internal/app/panel/domain/storage"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/modules/sslmanager/agent"
	"errors"
	"fmt"
	"strings"

	serverAgent "backend/internal/pkg/agent"
	"backend/internal/pkg/certificate"
	"backend/internal/pkg/logger"

	"github.com/r2dtools/agentintegration"
)

var ErrServerNotFound = errors.New("server not found")

type CertificateService struct {
	serverStorage         serverStorage.ServerStorage
	domainSettingsStorage domainStorage.DomainSettingStorage
	logger                logger.Logger
}

func (s CertificateService) IssueCertificate(request IssueCertificateRequest) (*dto.DomainCertificate, error) {
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return nil, err
	}

	emailSetting, err := s.domainSettingsStorage.FindByDomain(request.DomainName, request.ServerGuid, "email")

	if err != nil {
		return nil, err
	}

	if emailSetting == nil {
		err = s.domainSettingsStorage.Create(request.DomainName, request.ServerGuid, "email", request.Email)
	} else {
		emailSetting.SettingValue = request.Email
		err = s.domainSettingsStorage.Save(emailSetting)
	}

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
		return nil, err
	}

	return domainFactory.CreateCertificate(cert), nil
}

func (s CertificateService) AssignCertificate(request AssignCertificateRequest) (*dto.DomainCertificate, error) {
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return nil, err
	}

	cert, err := cAgent.AssignCertificateToDomain(&agentintegration.CertificateAssignRequestData{
		ServerName:  request.DomainName,
		WebServer:   request.WebServer,
		CertName:    request.CertName,
		StorageType: request.Storage,
	})

	if err != nil {
		return nil, err
	}

	return domainFactory.CreateCertificate(cert), nil
}

func (s CertificateService) UploadCertificate(
	guid string,
	requestData agentintegration.CertificateUploadRequestData,
) (*dto.DomainCertificate, error) {
	cAgent, err := s.getCertificateAgent(guid)

	if err != nil {
		return nil, err
	}

	cert, err := cAgent.Upload(&requestData)

	if err != nil {
		return nil, err
	}

	return domainFactory.CreateCertificate(cert), nil
}

func (s CertificateService) UploadCertificateToStorage(request UploadCertificateToStorageRequest) (*agentintegration.Certificate, error) {
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

func (s CertificateService) DownloadCertificateFromStorage(request DownloadCertificateRequest) (*agentintegration.CertificateDownloadResponseData, error) {
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return nil, err
	}

	requestData := agentintegration.CertificateDownloadRequestData{
		CertName:    request.CertName,
		StorageType: request.Storage,
	}

	return cAgent.DownloadtStorageCertificate(requestData)
}

func (s CertificateService) GetStorageCertificates(request CertificatesRequest) ([]StorageCertificateItem, error) {
	cAgent, err := s.getCertificateAgent(request.Guid)

	if err != nil {
		return nil, err
	}

	certsMap, err := cAgent.GetStorageCertificates()

	if err != nil {
		return nil, err
	}

	result := []StorageCertificateItem{}

	for name, cert := range certsMap {
		parts := strings.Split(name, "__")

		if len(parts) != 2 {
			return nil, errors.New("invalid certificate data")
		}

		result = append(result, StorageCertificateItem{
			Storage:     parts[0],
			CertName:    parts[1],
			Certificate: createCertificate(cert),
		})
	}

	return result, nil
}

func (s CertificateService) RemoveCertificateFromStorage(request RemoveCertificateFromStorageRequest) error {
	cAgent, err := s.getCertificateAgent(request.ServerGuid)

	if err != nil {
		return err
	}

	requestData := agentintegration.CertificateRemoveRequestData{
		CertName:    request.CertName,
		StorageType: request.Storage,
	}

	return cAgent.RemoveCertificateFromStorage(requestData)
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

	return cAgent.ChangeCommonDirStatus(agentintegration.CommonDirChangeStatusRequestData{
		WebServer:  request.WebServer,
		ServerName: request.DomainName,
		Status:     request.Status,
	})
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

func NewCertificateService(
	serverStorage serverStorage.ServerStorage,
	domainSettingStorage domainStorage.DomainSettingStorage,
	logger logger.Logger,
) CertificateService {
	return CertificateService{
		serverStorage:         serverStorage,
		domainSettingsStorage: domainSettingStorage,
		logger:                logger,
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
