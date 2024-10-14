package certificates

import (
	"backend/certificate"
	"backend/controllers"
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/r2dtools/agentintegration"
)

const (
	PEM_SUFFIX = ".pem"
	CRT_SUFFIX = ".crt"
)

// ModuleController controller for certificates module
type ModuleController struct {
	controllers.BaseController
}

func (ctrl *ModuleController) issueCertificate(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	var certData agentintegration.CertificateIssueRequestData
	if err = c.ShouldBindJSON(&certData); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	certData.ChallengeType = "http"
	certificate, err := agent.issue(&certData)

	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"certificate": certificate})
}

func (ctrl *ModuleController) assignCertificate(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	serverName := c.PostForm("ServerName")
	webServer := c.PostForm("WebServer")
	certName := c.PostForm("CertName")

	if serverName == "" || webServer == "" || certName == "" {
		ctrl.AbortWithInternalServerError(c, errors.New("invalid request data"))
		return
	}

	var requestData agentintegration.CertificateAssignRequestData
	requestData.ServerName = serverName
	requestData.WebServer = webServer
	requestData.CertName = certName
	certificate, err := agent.assignCertificateToDomain(&requestData)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"certificate": certificate})
}

func (ctrl *ModuleController) uploadCertificate(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	serverName := c.PostForm("ServerName")
	webServer := c.PostForm("WebServer")

	if serverName == "" || webServer == "" {
		ctrl.AbortWithInternalServerError(c, errors.New("invalid request data"))
		return
	}

	pemFileBytes, err := ctrl.getPemCertificateFromRequest(c)
	if err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	var requestData agentintegration.CertificateUploadRequestData
	requestData.PemCertificate = string(pemFileBytes)
	requestData.ServerName = serverName
	requestData.WebServer = webServer
	certificate, err := agent.upload(&requestData)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"certificate": certificate})
}

func (ctrl *ModuleController) uploadCertificateToStorage(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	certName := c.PostForm("CertName")

	if certName == "" {
		ctrl.AbortWithBadRequestError(c, errors.New("certificate name is missed"))
		return
	}
	certName = strings.TrimSuffix(certName, PEM_SUFFIX)
	certName = strings.TrimSuffix(certName, CRT_SUFFIX)

	pemFileBytes, err := ctrl.getPemCertificateFromRequest(c)
	if err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	certificate, err := ctrl.uploadCert(agent, string(pemFileBytes), certName)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"certificate": certificate})
}

func (ctrl *ModuleController) addSelfSignedCertificateToStorage(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}
	selfSignedCertData := c.PostForm("SelfSignedCert")
	if selfSignedCertData == "" {
		ctrl.AbortWithBadRequestError(c, errors.New("certificate data is missed"))
		return
	}

	certRequest := certificate.SelfSignedCertificateReuest{}
	if err = json.Unmarshal([]byte(selfSignedCertData), &certRequest); err != nil {
		ctrl.AbortWithInternalServerError(c, fmt.Errorf("invalid certificate data: %v", err))
		return
	}
	certPem, err := certificate.CreateSelfSignedCertificate(&certRequest)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, fmt.Errorf("could not generate self-signed certificate: %v", err))
		return
	}
	certificate, err := ctrl.uploadCert(agent, certPem, certRequest.CertName)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"certificate": certificate})
}

func (ctrl *ModuleController) removeCertificateFromStorage(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	data := struct{ CertName string }{}
	if err = c.ShouldBindJSON(&data); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	if data.CertName == "" {
		ctrl.AbortWithBadRequestError(c, errors.New("certificate name is missed"))
		return
	}

	if err = agent.removeCertificatefromStorage(data.CertName); err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}
}

func (ctrl *ModuleController) downloadCertificateFromStorage(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	data := struct{ CertName string }{}
	if err = c.ShouldBindJSON(&data); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	if data.CertName == "" {
		ctrl.AbortWithBadRequestError(c, errors.New("certificate name is missed"))
		return
	}

	certData, err := agent.downloadtStorageCertificate(data.CertName)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{
		"certFileName": certData.CertFileName,
		"certContent":  certData.CertContent,
	})
}

func (ctrl *ModuleController) getStorageCertNameList(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}
	certNameList, err := agent.getStorageCertNameList()
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"certNameList": certNameList})
}

func (ctrl *ModuleController) getStorageCertificateData(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.getAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	certName, ok := c.GetQuery("certName")
	if !ok {
		ctrl.AbortWithBadRequestError(c, errors.New("certificate name is not provided"))
		return
	}

	certificate, err := agent.getStorageCertificate(certName)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"certificate": certificate})
}

func (ctrl *ModuleController) getAgent(c *gin.Context, user *models.User) (*agent, error) {
	serverAgent, err := ctrl.GetServerAgent(c, user)

	if err != nil {
		return nil, err
	}

	return &agent{serverAgent}, nil
}

func (ctrl *ModuleController) getPemCertificateFromRequest(c *gin.Context) ([]byte, error) {
	c.Request.ParseMultipartForm(10)
	pemFile, _, err := c.Request.FormFile("PemCertificate")
	if err != nil {
		return nil, err
	}
	defer pemFile.Close()

	pemFileBytes, err := io.ReadAll(pemFile)
	if err != nil {
		return nil, err
	}

	return pemFileBytes, nil
}

func (ctrl *ModuleController) uploadCert(agent *agent, certPem, certName string) (*agentintegration.Certificate, error) {
	var requestData agentintegration.CertificateUploadRequestData
	requestData.PemCertificate = certPem
	requestData.CertName = certName

	return agent.uploadPemCertificateToStorage(&requestData)
}
