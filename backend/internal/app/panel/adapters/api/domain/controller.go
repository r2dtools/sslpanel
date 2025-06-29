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

		var request domainService.DomainRequest

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

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
				c.AbortWithError(http.StatusInternalServerError, err) // nolint:errcheck
			}

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

		var request domainService.DomainConfigRequest

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

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
				c.AbortWithError(http.StatusInternalServerError, err) // nolint:errcheck
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
		request := domainService.DomainSettingsRequest{
			DomainName: domainName,
			ServerGuid: guid,
		}
		settings, err := appDomainService.FindDomainSettings(request)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err) // nolint:errcheck

			return
		}

		settingsMap := map[string]string{}

		for _, setting := range settings {
			settingsMap[setting.SettingName] = setting.SettingValue
		}

		c.JSON(http.StatusOK, gin.H{"settings": settingsMap})
	}
}

func CreateChangeDomainSettingHandler(cAuth auth.Auth, appDomainService domainService.DomainService) func(c *gin.Context) {
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

		var request domainService.ChangeDomainSettingRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint:errcheck

			return
		}

		request.DomainName = domainName
		request.ServerGuid = guid

		err = appDomainService.ChangeDomainSettings(request)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

			return
		}
	}
}
