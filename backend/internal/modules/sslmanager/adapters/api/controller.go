package adapters

import (
	"backend/internal/app/panel/adapters/api/auth"
	"backend/internal/modules/sslmanager/service"
	"encoding/base64"
	"errors"
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		domainName := c.Param("domainName")

		if domainName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name")) // nolint:errcheck

			return
		}

		decodedDomainName, err := base64.RawStdEncoding.DecodeString(domainName)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name")) // nolint:errcheck

			return
		}

		domainName = string(decodedDomainName)

		var request service.CertificateIssueRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		request.ServerGuid = guid
		request.DomainName = domainName

		cert, err := certService.IssueCertificate(request)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		domainName := c.Param("domainName")

		if domainName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name")) // nolint:errcheck

			return
		}

		decodedDomainName, err := base64.RawStdEncoding.DecodeString(domainName)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name")) // nolint:errcheck

			return
		}

		domainName = string(decodedDomainName)

		var request service.CommonDirStatusRequest

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		request.ServerGuid = guid
		request.DomainName = domainName
		response, err := certService.GetCommonDirStatus(request)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		domainName := c.Param("domainName")

		if domainName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name")) // nolint:errcheck

			return
		}

		decodedDomainName, err := base64.RawStdEncoding.DecodeString(domainName)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name")) // nolint:errcheck

			return
		}

		domainName = string(decodedDomainName)

		var request service.ChangeCommonDirStatusRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		request.DomainName = domainName
		request.ServerGuid = guid
		err = certService.ChangeCommonDirStatus(request)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		domainName := c.Param("domainName")

		if domainName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name")) // nolint:errcheck

			return
		}

		decodedDomainName, err := base64.RawStdEncoding.DecodeString(domainName)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name")) // nolint:errcheck

			return
		}

		domainName = string(decodedDomainName)

		var request service.AssignCertificateRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		request.ServerGuid = guid
		request.DomainName = domainName
		_, err = certService.AssignCertificate(request)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
		}
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		serverName := c.PostForm("ServerName")
		webServer := c.PostForm("WebServer")

		if serverName == "" || webServer == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid request data")) // nolint:errcheck

			return
		}

		pemFileBytes, err := getPemCertificateFromRequest(c)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		var requestData agentintegration.CertificateUploadRequestData

		requestData.PemCertificate = string(pemFileBytes)
		requestData.ServerName = serverName
		requestData.WebServer = webServer

		cert, err := certService.UploadCertificate(guid, requestData)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		certName := c.PostForm("name")

		if certName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("certificate name is missed")) // nolint:errcheck

			return
		}

		certName = strings.TrimSuffix(certName, pemSuffix)
		certName = strings.TrimSuffix(certName, crtSuffix)

		pemFileBytes, err := getPemCertificateFromRequest(c)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		request := service.CertificateUploadToStorageRequest{
			ServerGuid:     guid,
			CertName:       certName,
			PemCertificate: string(pemFileBytes),
		}
		_, err = certService.UploadCertificateToStorage(request)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
		}
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		requestData := struct {
			CertName string `json:"name"`
		}{}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		if requestData.CertName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("certificate name is missed")) // nolint:errcheck
			return
		}

		certData, err := certService.DownloadCertificateFromStorage(guid, requestData.CertName)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"name":    certData.CertFileName,
			"content": certData.CertContent,
		})
	}
}

func CreateGetStorageCertificatesHandler(cAuth auth.Auth, certService service.CertificateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		guid := c.Param("serverId")

		if guid == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		request := service.CertificatesRequest{Guid: guid}
		certsMap, err := certService.GetStorageCertificates(request)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"certificates": certsMap})
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		certName, ok := c.GetQuery("certName")

		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid certificate name")) // nolint:errcheck

			return
		}

		certificate, err := certService.GetStorageCertificate(guid, certName)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		requestData := struct{ CertName string }{}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		if requestData.CertName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("certificate name is missed")) // nolint:errcheck

			return
		}

		err := certService.RemoveCertificateFromStorage(guid, requestData.CertName)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server GUID")) // nolint:errcheck

			return
		}

		request := service.SelfSignedCertificateRequest{}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		if request.CertName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("certificate name is missed")) // nolint:errcheck

			return
		}

		if request.CommonName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("common name is missed")) // nolint:errcheck

			return
		}

		request.ServerGuid = guid
		_, err := certService.CreateSelfSignCertificate(request)

		if err != nil {
			if errors.Is(err, service.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}

			return
		}
	}
}

func getPemCertificateFromRequest(c *gin.Context) ([]byte, error) {
	err := c.Request.ParseMultipartForm(10)

	if err != nil {
		return nil, err
	}

	pemFile, _, err := c.Request.FormFile("file")

	if err != nil {
		return nil, err
	}

	defer pemFile.Close() // nolint:errcheck

	return io.ReadAll(pemFile)
}
