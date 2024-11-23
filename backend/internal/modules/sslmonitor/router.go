package sslmonitor

import (
	"backend/internal/app/panel/adapters/api/auth"
	monitorApi "backend/internal/modules/sslmonitor/adapters/api"
	"backend/internal/modules/sslmonitor/service"
	domainStorage "backend/internal/modules/sslmonitor/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(
	group *gin.RouterGroup,
	cAuth auth.Auth,
	db *gorm.DB,
) {
	appDomainStorage := domainStorage.NewDomainStorage(db)
	appMonitorService := service.NewMonitorService(appDomainStorage)

	group.GET("/sites", monitorApi.CreateGetUserDomainsHandler(cAuth, appMonitorService))
	group.DELETE("/sites/:siteId", monitorApi.CreateRemoveDomainHandler(cAuth, appMonitorService))
	group.POST("/sites", monitorApi.CreateAddDomainHandler(cAuth, appMonitorService))
	group.POST("/sites/:siteId/refresh", monitorApi.CreateRefreshDomainHandler(cAuth, appMonitorService))
}
