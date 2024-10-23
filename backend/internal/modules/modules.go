package modules

import (
	"backend/internal/app/panel/adapters/api/auth"
	serverStorage "backend/internal/app/panel/server/storage"
	userStorage "backend/internal/app/panel/user/storage"
	sslManagerModule "backend/internal/modules/sslmanager"
	sslMonitorModule "backend/internal/modules/sslmonitor"
	sslMonitorCli "backend/internal/modules/sslmonitor/adapters/cli"
	"backend/internal/pkg/logger"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

func InitModulesRouter(
	group *gin.RouterGroup,
	db *gorm.DB,
	authMiddleware *jwt.GinJWTMiddleware,
	cAuth auth.Auth,
	appServerStorage serverStorage.ServerStorage,
) {
	certificatesGroup := group.Group("certificates")
	{
		certificatesGroup.Use(authMiddleware.MiddlewareFunc())
		sslManagerModule.InitRouter(certificatesGroup, cAuth, appServerStorage)
	}

	sslMonitorGroup := group.Group("certificate-monitor")
	{
		sslMonitorGroup.Use(authMiddleware.MiddlewareFunc())
		sslMonitorModule.InitRouter(sslMonitorGroup, cAuth, db)
	}
}

func InitModulesCli(
	rootCmd *cobra.Command,
	db *gorm.DB,
	appUserStorage userStorage.UserStorage,
	logger logger.Logger,
) {
	sslMonitorCli.InitCmd(rootCmd, db, appUserStorage, logger)
}
