package certificatemonitor

import (
	"github.com/gin-gonic/gin"
)

// InitRouter registers Certificate Monitor module routes
func InitRouter(group *gin.RouterGroup) {
	ctrl := new(ModuleController)
	group.GET("/sites", ctrl.getSites)
	group.DELETE("/sites/:siteId", ctrl.removeSite)
	group.POST("/sites", ctrl.createSite)
	group.POST("/sites/:siteId/refresh", ctrl.refreshSite)
}
