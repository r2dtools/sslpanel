package controllers

import (
	"backend/models"
	"backend/remote"
	"backend/tcp"
	"backend/utils"
	"errors"

	"github.com/gin-gonic/gin"
)

// ServerController is acontroller to manage servers
type ServerController struct {
	BaseController
}

// GetServers return all registered servers for the currently authenticated user
func (ctrl *ServerController) GetServers(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	servers, err := models.GetAccountServers(user.AccountID)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"servers": servers})
}

// GetServer return all registered servers for the currently authenticated user
func (ctrl *ServerController) GetServer(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	server, err := ctrl.GetServerByID(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"server": server})
}

// RemoveServer remove the server of a current user by ID
func (ctrl *ServerController) RemoveServer(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	server, err := ctrl.GetServerByID(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	if err = server.Remove(); err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}
}

// AddServer adds new unregistered server
func (ctrl *ServerController) AddServer(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	var server models.Server
	var err error
	if err = c.ShouldBindJSON(&server); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	server.AccountID = user.AccountID
	err = models.SaveServer(&server)

	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}
}

// SaveServer saves changes for server
func (ctrl *ServerController) SaveServer(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	server, err := ctrl.GetServerByID(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	if err = c.ShouldBindJSON(server); err != nil {
		ctrl.AbortWithBadRequestError(c, err)
		return
	}

	err = models.SaveServer(server)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"server": server})
}

// RefreshServer registers new server
func (ctrl *ServerController) RefreshServer(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	server, err := ctrl.GetServerByID(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	agent := &remote.Agent{
		Server: server,
	}
	refData, err := agent.Refresh()
	var connErr *tcp.ConnectionError

	if err != nil {
		// If connection failed with server agent then change agent status to "inactive"
		if errors.As(err, &connErr) {
			defer (func() {
				server.IsActive = 0
				models.SaveServer(server)
			})()
		}

		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	server.OsCode = refData.Platform
	server.OsVersion = refData.PlatformVersion
	server.AgentVersion = refData.AgentVersion
	server.IsRegistered = 1
	server.IsActive = 1
	if err = models.SaveServer(server); err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"server": server, "serverData": refData})
}

// GetServerVhosts loads server virtual hosts
func (ctrl *ServerController) GetServerVhosts(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	agent, err := ctrl.GetServerAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	vhosts, err := agent.GetVhosts()
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"vhosts": utils.FilterVhosts(vhosts)})
}

// GetVhostCertificate loads vhost certificate
func (ctrl *ServerController) GetVhostCertificate(c *gin.Context) {
	user := ctrl.GetCurrentUser(c)
	if user == nil {
		return
	}

	vhostName, ok := c.GetQuery("vhostName")
	if !ok {
		ctrl.AbortWithBadRequestError(c, errors.New("vhostName parameter is missed"))
		return
	}

	agent, err := ctrl.GetServerAgent(c, user)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	certificate, err := agent.GetVhostCertificate(vhostName)
	if err != nil {
		ctrl.AbortWithInternalServerError(c, err)
		return
	}

	ctrl.SuccessJSON(c, gin.H{"certificate": certificate})
}
