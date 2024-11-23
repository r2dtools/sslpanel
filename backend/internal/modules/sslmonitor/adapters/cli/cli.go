package cli

import (
	userStorage "backend/internal/app/panel/user/storage"
	monitorService "backend/internal/modules/sslmonitor/service"
	domainStorage "backend/internal/modules/sslmonitor/storage"
	"backend/internal/pkg/logger"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
	Use:   "sslmonitor",
	Short: "SSL Monitor CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func InitCmd(
	rCmd *cobra.Command,
	db *gorm.DB,
	appUserStorage userStorage.UserStorage,
	logger logger.Logger,
) {
	appDomainStorage := domainStorage.NewDomainStorage(db)
	appMonitorService := monitorService.NewMonitorService(appDomainStorage)

	rootCmd.AddCommand(getRefreshCmd(appMonitorService, appDomainStorage, logger))
	rootCmd.AddCommand(getNotifyCmd(appUserStorage, appMonitorService, logger))

	rCmd.AddCommand(rootCmd)
}
