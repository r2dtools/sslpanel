package adapters

import (
	"backend/internal/app/panel/adapters/api/auth"
	"backend/internal/modules/sslmanager/service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		var certData agentintegration.CertificateIssueRequestData

		if err := c.ShouldBindJSON(&certData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		cert, err := certService.IssueCertificate(serverID, certData)

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

func CreateAssignCertificateHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		var requestData agentintegration.CertificateAssignRequestData

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		cert, err := certService.AssignCertificate(serverID, requestData)

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

		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

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

		cert, err := certService.UploadCertificate(serverID, requestData)

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

		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

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

		certificate, err := certService.UploadCertificateToStorage(serverID, requestData)

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

		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		requestData := struct{ CertName string }{}

		if err = c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		if requestData.CertName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("certificate name is missed"))
			return
		}

		certData, err := certService.DownloadCertificateFromStorage(serverID, requestData.CertName)

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

		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		certNameList, err := certService.GetStorageCertNameList(serverID)

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

		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		certName, ok := c.GetQuery("certName")

		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid certificate name"))

			return
		}

		certificate, err := certService.GetStorageCertificate(serverID, certName)

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

		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		requestData := struct{ CertName string }{}

		if err = c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		if requestData.CertName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("certificate name is missed"))

			return
		}

		err = certService.RemoveCertificateFromStorage(serverID, requestData.CertName)

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

		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		selfSignedCertData := c.PostForm("SelfSignedCert")

		if selfSignedCertData == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid certificate data"))

			return
		}

		requestData := service.SelfSignedCertificateRequest{}

		if err = json.Unmarshal([]byte(selfSignedCertData), &requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid certificate data: %v", err))

			return
		}

		certificate, err := certService.CreateSelfSignCertificate(serverID, requestData)

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
