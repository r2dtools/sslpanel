package domain

import (
	"backend/internal/app/panel/adapters/api/auth"
	domainService "backend/internal/app/panel/domain/service"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGetDomainHandler(cAuth auth.Auth, appDomainService domainService.DomainService) func(c *gin.Context) {
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

		domainName := c.Param("domainName")

		if domainName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name"))

			return
		}

		decodedDomainName, err := base64.RawStdEncoding.DecodeString(domainName)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name"))

			return
		}

		var request domainService.DomainRequest

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		request.ServerGuid = guid
		request.DomainName = string(decodedDomainName)
		domain, err := appDomainService.GetDomain(request)

		if err != nil {
			if errors.Is(err, domainService.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else if errors.Is(err, domainService.ErrDomainNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else if errors.Is(err, domainService.ErrAgentConnection) {
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		if domain == nil {
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"domain": domain})
	}
}

func CreateGetDomainConfigHandler(cAuth auth.Auth, appService domainService.DomainService) func(c *gin.Context) {
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

		domainName := c.Param("domainName")

		if domainName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name"))

			return
		}

		decodedDomainName, err := base64.RawStdEncoding.DecodeString(domainName)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name"))

			return
		}

		var request domainService.DomainConfigRequest

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		request.ServerGuid = guid
		request.DomainName = string(decodedDomainName)
		config, err := appService.GetDomainConfig(request)
		var errAgentCommon domainService.ErrAgentCommon

		if err != nil {
			if errors.Is(err, domainService.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else if errors.As(err, &errAgentCommon) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"content": config})
	}
}

func CreateFindDomainSettingsHandler(cAuth auth.Auth, appDomainService domainService.DomainService) func(c *gin.Context) {
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

		domainName := c.Param("domainName")

		if domainName == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name"))

			return
		}

		decodedDomainName, err := base64.RawStdEncoding.DecodeString(domainName)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain name"))

			return
		}

		domainName = string(decodedDomainName)
		settings, err := appDomainService.FindDomainSettings(domainName, guid)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"settings": settings})
	}
}
