package certificates

import (
	"github.com/gin-gonic/gin"
)

// InitRouter registers certificates module routes
func InitRouter(group *gin.RouterGroup) {
	ctrl := new(ModuleController)
	group.POST("/:serverId/issue/:serverName", ctrl.issueCertificate)
	group.POST("/:serverId/upload/:serverName", ctrl.uploadCertificate)
	group.POST("/:serverId/domain/assign", ctrl.assignCertificate)
	group.GET("/:serverId/storage/cert-name-list", ctrl.getStorageCertNameList)
	group.GET("/:serverId/storage/cert-data", ctrl.getStorageCertificateData)
	group.POST("/:serverId/storage/upload", ctrl.uploadCertificateToStorage)
	group.POST("/:serverId/storage/remove", ctrl.removeCertificateFromStorage)
	group.POST("/:serverId/storage/add-self-signed", ctrl.addSelfSignedCertificateToStorage)
	group.POST("/:serverId/storage/download", ctrl.downloadCertificateFromStorage)
}
