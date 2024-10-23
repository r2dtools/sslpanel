package sslmanager

import (
	"backend/internal/app/panel/adapters/api/auth"
	serverStorage "backend/internal/app/panel/server/storage"
	certApi "backend/internal/modules/sslmanager/adapters/api"
	"backend/internal/modules/sslmanager/service"

	"github.com/gin-gonic/gin"
)

func InitRouter(
	group *gin.RouterGroup,
	cAuth auth.Auth,
	appServerStorage serverStorage.ServerStorage,
) {
	appCertificateService := service.NewCertificateService(appServerStorage)

	group.POST("/:serverId/issue/:serverName", certApi.CreateIssueCertificateHandler(cAuth, appCertificateService))
	group.POST("/:serverId/domain/assign", certApi.CreateAssignCertificateHandler(cAuth, appCertificateService))
	group.POST("/:serverId/upload/:serverName", certApi.CreateUploadCertificateHandler(cAuth, appCertificateService))
	group.POST("/:serverId/storage/upload", certApi.CreateUploadCertificateToStorageHandler(cAuth, appCertificateService))
	group.POST("/:serverId/storage/download", certApi.CreateDownloadCertificateFromStorageHandler(cAuth, appCertificateService))
	group.GET("/:serverId/storage/cert-name-list", certApi.CreateGetStorageCertNameListHandler(cAuth, appCertificateService))
	group.GET("/:serverId/storage/cert-data", certApi.CreateGetStorageCertNameListHandler(cAuth, appCertificateService))
	group.POST("/:serverId/storage/remove", certApi.CreateRemoveCertificateFromStorageHandler(cAuth, appCertificateService))
	group.POST("/:serverId/storage/add-self-signed", certApi.CreateAddSelfSignCertificateToStorageHandler(cAuth, appCertificateService))
}
