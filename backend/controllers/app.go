package controllers

import (
	"backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BaseController  is acontroller to manage account settings
type AppController struct {
	BaseController
}

// ChangePassword changes user`s password
func (ctrl AppController) GetData(c *gin.Context) {
	if !ctrl.CheckAuth(c) {
		return
	}

	aConfig := config.GetConfig()
	data := struct {
		AgentVersion          string `json:"agentVersion"`
		AgentInstallerVersion string `json:"agentInstallerVersion"`
	}{AgentVersion: aConfig.AgentVersion, AgentInstallerVersion: aConfig.InstallerVersion}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
