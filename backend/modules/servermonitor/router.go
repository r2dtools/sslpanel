package servermonitor

import (
	"github.com/gin-gonic/gin"
)

// InitRouter registers Server Monitor module routes
func InitRouter(group *gin.RouterGroup) {
	ctrl := new(ModuleController)
	group.GET("/:serverId/category-data", ctrl.getCategoryData)
}
