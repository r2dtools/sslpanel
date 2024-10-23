package api

import (
	"backend/internal/app/panel/adapters/api/auth"
	"backend/internal/modules/sslmonitor/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGetUserDomainsHandler(cAuth auth.Auth, monitorService service.MonitorService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		domains, err := monitorService.FindUserDomains(user.ID)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"sites": domains})
	}
}

func CreateRemoveDomainHandler(cAuth auth.Auth, monitorService service.MonitorService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		id, err := strconv.Atoi(c.Param("siteId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain id"))

			return
		}

		err = monitorService.RemoveDomain(id)

		if err != nil {
			if errors.Is(err, service.ErrDomainNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func CreateAddDomainHandler(cAuth auth.Auth, monitorService service.MonitorService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		var requestData service.AddDomainRequest

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err := monitorService.AddDomain(requestData, user.ID)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

	}
}

func CreateRefreshDomainHandler(cAuth auth.Auth, monitorService service.MonitorService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		id, err := strconv.Atoi(c.Param("siteId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid domain id"))

			return
		}

		domain, err := monitorService.RefreshDomain(id)

		if err != nil {
			if errors.Is(err, service.ErrDomainNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}

		c.JSON(http.StatusOK, gin.H{"site": domain})
	}
}
