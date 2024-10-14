package modules

import (
	certificatemonitormodule "backend/modules/certificatemonitor"
	certificatemonitormodulecdm "backend/modules/certificatemonitor/cmd"
	certificatesmodule "backend/modules/certificates"
	servermonitormodule "backend/modules/servermonitor"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// InitModulesRouter register modules routes
func InitModulesRouter(group *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	certificatesGroup := group.Group("certificates")
	{
		certificatesGroup.Use(authMiddleware.MiddlewareFunc())
		certificatesmodule.InitRouter(certificatesGroup)
	}
	certificateMonitorGroup := group.Group("certificate-monitor")
	{
		certificateMonitorGroup.Use(authMiddleware.MiddlewareFunc())
		certificatemonitormodule.InitRouter(certificateMonitorGroup)
	}
	serverMonitorGroup := group.Group("server-monitor")
	{
		serverMonitorGroup.Use(authMiddleware.MiddlewareFunc())
		servermonitormodule.InitRouter(serverMonitorGroup)
	}
}

// InitModulesCli registers modules cli
func InitModulesCli(rootCmd *cobra.Command) {
	certificatemonitormodulecdm.InitCmd(rootCmd)
}
