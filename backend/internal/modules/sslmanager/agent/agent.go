package agent

import (
	"fmt"

	serverAgent "backend/internal/pkg/agent"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agentintegration"
)

const (
	HttpChallenge = "http"
)

type CertificateAgent struct {
	serverAgent *serverAgent.Agent
}

func (a *CertificateAgent) Issue(certData agentintegration.CertificateIssueRequestData) (*agentintegration.Certificate, error) {
	data, err := a.serverAgent.Request("certificates.issue", certData)

	if err != nil {
		return nil, err
	}

	return getCertificate(data)
}

func (a *CertificateAgent) Upload(certData *agentintegration.CertificateUploadRequestData) (*agentintegration.Certificate, error) {
	data, err := a.serverAgent.Request("certificates.upload", certData)

	if err != nil {
		return nil, err
	}

	return getCertificate(data)
}

func (a *CertificateAgent) AssignCertificateToDomain(certData *agentintegration.CertificateAssignRequestData) (*agentintegration.Certificate, error) {
	data, err := a.serverAgent.Request("certificates.domainassign", certData)

	if err != nil {
		return nil, err
	}

	return getCertificate(data)
}

func (a *CertificateAgent) UploadPemCertificateToStorage(certData *agentintegration.CertificateUploadRequestData) (*agentintegration.Certificate, error) {
	data, err := a.serverAgent.Request("certificates.storagecertupload", certData)

	if err != nil {
		return nil, err
	}

	return getCertificate(data)
}

func (a *CertificateAgent) RemoveCertificateFromStorage(request agentintegration.CertificateRemoveRequestData) error {
	_, err := a.serverAgent.Request("certificates.storagecertremove", request)

	return err
}

func (a *CertificateAgent) GetStorageCertificates() (map[string]*agentintegration.Certificate, error) {
	certificates := map[string]*agentintegration.Certificate{}
	data, err := a.serverAgent.Request("certificates.storagecertificates", nil)

	if err != nil {
		return nil, err
	}

	if data == nil {
		return certificates, nil
	}

	var response agentintegration.CertificatesResponseData
	err = mapstructure.Decode(data, &response)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate list data: %v", err)
	}

	return response.Certificates, nil
}

func (a *CertificateAgent) DownloadtStorageCertificate(request agentintegration.CertificateDownloadRequestData) (*agentintegration.CertificateDownloadResponseData, error) {
	data, err := a.serverAgent.Request("certificates.storagecertdownload", request)

	if err != nil {
		return nil, err
	}

	var certData agentintegration.CertificateDownloadResponseData
	err = mapstructure.Decode(data, &certData)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate name list data: %v", err)
	}

	return &certData, nil
}

func (a *CertificateAgent) GetCommonDirStatus(request agentintegration.CommonDirStatusRequestData) (agentintegration.CommonDirStatusResponseData, error) {
	var responsse agentintegration.CommonDirStatusResponseData

	data, err := a.serverAgent.Request("certificates.commondirstatus", request)

	if err != nil {
		return responsse, err
	}

	err = mapstructure.Decode(data, &responsse)

	if err != nil {
		return responsse, fmt.Errorf("invalid common dir status data: %v", err)
	}

	return responsse, nil
}

func (a *CertificateAgent) ChangeCommonDirStatus(request agentintegration.CommonDirChangeStatusRequestData) error {
	_, err := a.serverAgent.Request("certificates.changecommondirstatus", request)

	return err
}

func getCertificate(responseData interface{}) (*agentintegration.Certificate, error) {
	if responseData == nil {
		return nil, nil
	}

	var certificate agentintegration.Certificate
	err := mapstructure.Decode(responseData, &certificate)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate data: %v", err)
	}

	return &certificate, nil
}

func NewCertificateAgent(serverAgent *serverAgent.Agent) *CertificateAgent {
	return &CertificateAgent{
		serverAgent: serverAgent,
	}
}
