package certificates

import (
	"backend/remote"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agentintegration"
)

type agent struct {
	serverAgent *remote.Agent
}

func (a *agent) issue(certData *agentintegration.CertificateIssueRequestData) (*agentintegration.Certificate, error) {
	data, err := a.serverAgent.Request("certificates.issue", certData)
	if err != nil {
		return nil, err
	}

	return getCertificate(data)
}

func (a *agent) upload(certData *agentintegration.CertificateUploadRequestData) (*agentintegration.Certificate, error) {
	data, err := a.serverAgent.Request("certificates.upload", certData)
	if err != nil {
		return nil, err
	}

	return getCertificate(data)
}

func (a *agent) assignCertificateToDomain(certData *agentintegration.CertificateAssignRequestData) (*agentintegration.Certificate, error) {
	data, err := a.serverAgent.Request("certificates.domainassign", certData)
	if err != nil {
		return nil, err
	}

	return getCertificate(data)
}

func (a *agent) uploadPemCertificateToStorage(certData *agentintegration.CertificateUploadRequestData) (*agentintegration.Certificate, error) {
	data, err := a.serverAgent.Request("certificates.storagecertupload", certData)
	if err != nil {
		return nil, err
	}

	return getCertificate(data)
}

func (a *agent) removeCertificatefromStorage(certName string) error {
	_, err := a.serverAgent.Request("certificates.storagecertremove", certName)
	if err != nil {
		return err
	}

	return nil
}

func (a *agent) getStorageCertNameList() ([]string, error) {
	data, err := a.serverAgent.Request("certificates.storagecertnamelist", nil)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return []string{}, nil
	}

	var certNameList agentintegration.StorageCertificateNameList
	err = mapstructure.Decode(data, &certNameList)
	if err != nil {
		return nil, fmt.Errorf("invalid certificate name list data: %v", err)
	}

	return certNameList.CertNameList, nil
}

func (a *agent) getStorageCertificate(certName string) (*agentintegration.Certificate, error) {
	data, err := a.serverAgent.Request("certificates.storagecertdata", certName)
	if err != nil {
		return nil, err
	}

	return getCertificate(data)
}

func (a *agent) downloadtStorageCertificate(certName string) (*agentintegration.CertificateDownloadResponseData, error) {
	data, err := a.serverAgent.Request("certificates.storagecertdownload", certName)
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
