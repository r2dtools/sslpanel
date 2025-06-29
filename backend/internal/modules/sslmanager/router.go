package sslmanager

import (
	"backend/internal/app/panel/adapters/api/auth"
	domainStorage "backend/internal/app/panel/domain/storage"
	serverStorage "backend/internal/app/panel/server/storage"
	certApi "backend/internal/modules/sslmanager/adapters/api"
	"backend/internal/modules/sslmanager/service"
	"backend/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

func InitRouter(
	group *gin.RouterGroup,
	cAuth auth.Auth,
	appServerStorage serverStorage.ServerStorage,
	appDomainSettingStorage domainStorage.DomainSettingStorage,
	logger logger.Logger,
) {
	appCertificateService := service.NewCertificateService(appServerStorage, appDomainSettingStorage, logger)

	group.POST("/:serverId/domain/:domainName/issue", certApi.CreateIssueCertificateHandler(cAuth, appCertificateService))
	group.POST("/:serverId/domain/:domainName/assign", certApi.CreateAssignCertificateHandler(cAuth, appCertificateService))
	group.GET("/:serverId/domain/:domainName/commondir-status", certApi.CreateGetCommonDirStatusHandler(cAuth, appCertificateService))
	group.POST("/:serverId/domain/:domainName/commondir-status", certApi.CreateChangeCommonDirStatusHandler(cAuth, appCertificateService))
	group.POST("/:serverId/upload/:serverName", certApi.CreateUploadCertificateHandler(cAuth, appCertificateService))
	group.POST("/:serverId/storage/upload", certApi.CreateUploadCertificateToStorageHandler(cAuth, appCertificateService))
	group.POST("/:serverId/storage/download", certApi.CreateDownloadCertificateFromStorageHandler(cAuth, appCertificateService))
	group.GET("/:serverId/storage/certificates", certApi.CreateGetStorageCertificatesHandler(cAuth, appCertificateService))
	group.POST("/:serverId/storage/remove", certApi.CreateRemoveCertificateFromStorageHandler(cAuth, appCertificateService))
	group.POST("/:serverId/storage/add-self-signed", certApi.CreateAddSelfSignCertificateToStorageHandler(cAuth, appCertificateService))
}
