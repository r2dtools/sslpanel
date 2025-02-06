package adapters

import (
	"backend/internal/app/panel/adapters/api/auth"
	"backend/internal/modules/sslmanager/service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/r2dtools/agentintegration"
)

const (
	pemSuffix = ".pem"
	crtSuffix = ".crt"
)

func CreateIssueCertificateHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		var certIssueRequest service.CertificateIssueRequest

		if err := c.ShouldBindJSON(&certIssueRequest); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		var errAgentCommon service.ErrAgentCommon
		cert, err := certService.IssueCertificate(guid, certIssueRequest)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else if errors.As(err, &errAgentCommon) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"certificate": cert})
	}
}

func CreateGetCommonDirStatusHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		var request service.CommonDirStatusRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		response, err := certService.GetCommonDirStatus(guid, request)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"commondir": response})
	}
}

func CreateChangeCommonDirStatusHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		var request service.CommonDirStatusChangeRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err := certService.ChangeCommonDirStatus(guid, request)

		var errAgentCommon service.ErrAgentCommon

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else if errors.As(err, &errAgentCommon) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func CreateAssignCertificateHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		var requestData agentintegration.CertificateAssignRequestData

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		cert, err := certService.AssignCertificate(guid, requestData)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"certificate": cert})
	}
}

func CreateUploadCertificateHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		serverName := c.PostForm("ServerName")
		webServer := c.PostForm("WebServer")

		if serverName == "" || webServer == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid request data"))

			return
		}

		pemFileBytes, err := getPemCertificateFromRequest(c)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		var requestData agentintegration.CertificateUploadRequestData

		requestData.PemCertificate = string(pemFileBytes)
		requestData.ServerName = serverName
		requestData.WebServer = webServer

		cert, err := certService.UploadCertificate(guid, requestData)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"certificate": cert})
	}
}

func CreateUploadCertificateToStorageHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		certName := c.PostForm("CertName")

		if certName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("certificate name is missed"))

			return
		}

		certName = strings.TrimSuffix(certName, pemSuffix)
		certName = strings.TrimSuffix(certName, crtSuffix)

		pemFileBytes, err := getPemCertificateFromRequest(c)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		var requestData agentintegration.CertificateUploadRequestData
		requestData.PemCertificate = string(pemFileBytes)
		requestData.CertName = certName

		certificate, err := certService.UploadCertificateToStorage(guid, requestData)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"certificate": certificate})
	}
}

func CreateDownloadCertificateFromStorageHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		requestData := struct{ CertName string }{}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		if requestData.CertName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("certificate name is missed"))
			return
		}

		certData, err := certService.DownloadCertificateFromStorage(guid, requestData.CertName)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"certFileName": certData.CertFileName,
			"certContent":  certData.CertContent,
		})
	}
}

func CreateGetStorageCertNameListHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		certNameList, err := certService.GetStorageCertNameList(guid)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"certNameList": certNameList})
	}
}

func CreateGetStorageCertificateHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		certName, ok := c.GetQuery("certName")

		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid certificate name"))

			return
		}

		certificate, err := certService.GetStorageCertificate(guid, certName)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"certificate": certificate})
	}
}

func CreateRemoveCertificateFromStorageHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		requestData := struct{ CertName string }{}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		if requestData.CertName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("certificate name is missed"))

			return
		}

		err := certService.RemoveCertificateFromStorage(guid, requestData.CertName)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}
	}
}

func CreateAddSelfSignCertificateToStorageHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID"))

			return
		}

		selfSignedCertData := c.PostForm("SelfSignedCert")

		if selfSignedCertData == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid certificate data"))

			return
		}

		requestData := service.SelfSignedCertificateRequest{}

		if err := json.Unmarshal([]byte(selfSignedCertData), &requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid certificate data: %v", err))

			return
		}

		certificate, err := certService.CreateSelfSignCertificate(guid, requestData)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"certificate": certificate})
	}
}

func getPemCertificateFromRequest(c *gin.Context) ([]byte, error) {
	c.Request.ParseMultipartForm(10)
	pemFile, _, err := c.Request.FormFile("PemCertificate")

	if err != nil {
		return nil, err
	}

	defer pemFile.Close()

	return io.ReadAll(pemFile)
}
