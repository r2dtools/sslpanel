package server

import (
	"backend/internal/app/panel/adapters/api/auth"
	serverService "backend/internal/app/panel/server/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
)

func CreateFindAccounServersHandler(cAuth auth.Auth, appServerService serverService.ServerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		servers, err := appServerService.FindAccountServers(user.AccountID)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"servers": servers})
	}
}

func CreateGetServerByGuidHandler(cAuth auth.Auth, appServerService serverService.ServerService) func(c *gin.Context) {
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

		server, err := appServerService.FindServerByGuid(guid)

		if err != nil {
			if errors.Is(err, serverService.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		if server == nil {
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"server": server})
	}
}

func CreateGetServerDetailsByGuidHandler(cAuth auth.Auth, appServerService serverService.ServerService) func(c *gin.Context) {
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

		server, err := appServerService.GetServerDetailsByGuid(guid)

		if err != nil {
			if errors.Is(err, serverService.ErrServerNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else if errors.Is(err, serverService.ErrAgentConnection) {
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		if server == nil {
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"server": server})
	}
}

func CreateRemoveServerHandler(cAuth auth.Auth, appServerService serverService.ServerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		id, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		err = appServerService.RemoveServer(id)

		if err != nil {
			if errors.Is(err, serverService.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}
	}
}

func CreateAddServerHandler(cAuth auth.Auth, appServerService serverService.ServerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := cAuth.GetCurrentUser(c)

		if user == nil {
			return
		}

		var request serverService.NewServerRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err := validator.Validate(request)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		request.AccountID = user.AccountID

		err = appServerService.AddServer(request)

		if err != nil {
			var errService serverService.ErrServerService

			if errors.As(err, &errService) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func CreateUpdateServerHandler(appServerService serverService.ServerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		var request serverService.UpdateServerRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		request.ID = serverID
		err = validator.Validate(request)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		err = appServerService.UpdateServer(request)

		if err != nil {
			var errService serverService.ErrServerService

			if errors.As(err, &errService) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func CreateGetVhostCertificateHandler(appServerService serverService.ServerService) func(c *gin.Context) {
	return func(c *gin.Context) {
		serverID, err := strconv.Atoi(c.Param("serverId"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid server ID"))

			return
		}

		vhostName, ok := c.GetQuery("vhostName")

		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid virtual host name"))

			return
		}

		certificate, err := appServerService.GetVhostCertificate(serverID, vhostName)

		if err != nil {
			if errors.Is(err, serverService.ErrServerNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"certificate": certificate})
	}
}
