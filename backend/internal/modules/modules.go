package modules

import (
	"backend/internal/app/panel/adapters/api/auth"
	domainStorage "backend/internal/app/panel/domain/storage"
	serverStorage "backend/internal/app/panel/server/storage"
	sslManagerModule "backend/internal/modules/sslmanager"
	"backend/internal/modules/sslmanager/autorenewal/logstorage"
	"backend/internal/pkg/logger"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModulesRouter(
	group *gin.RouterGroup,
	db *gorm.DB,
	authMiddleware *jwt.GinJWTMiddleware,
	cAuth auth.Auth,
	appServerStorage serverStorage.ServerStorage,
	appDomainSettingStorage domainStorage.DomainSettingStorage,
	certRenewalLogStorage logstorage.RenewalLogStorage,
	logger logger.Logger,
) {
	certificatesGroup := group.Group("certificates")
	{
		certificatesGroup.Use(authMiddleware.MiddlewareFunc())
		sslManagerModule.InitRouter(
			certificatesGroup,
			cAuth,
			appServerStorage,
			appDomainSettingStorage,
			certRenewalLogStorage,
			logger,
		)
	}
}
