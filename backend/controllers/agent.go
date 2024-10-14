package controllers

import (
	"backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AgentController  is a controller to manage account settings
type AgentController struct {
	BaseController
}

func (ctrl AgentController) LatestVersion(c *gin.Context) {
	aConfig := config.GetConfig()
	c.JSON(http.StatusOK, gin.H{"version": aConfig.AgentVersion})
}
