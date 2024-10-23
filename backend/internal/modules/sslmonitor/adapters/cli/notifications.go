package cli

import (
	userStorage "backend/internal/app/panel/user/storage"
	monitorService "backend/internal/modules/sslmonitor/service"
	"backend/internal/pkg/logger"
	"fmt"

	"github.com/spf13/cobra"
)

func getNotifyCmd(
	appUserStorage userStorage.UserStorage,
	appMonitorService monitorService.MonitorService,
	logger logger.Logger,
) *cobra.Command {
	var notificationsCmd = &cobra.Command{
		Use:   "notify",
		Short: "Notify user about SSL expiration",
		RunE: func(cmd *cobra.Command, args []string) error {
			users, err := appUserStorage.FindAll()

			if err != nil {
				logger.Debug(fmt.Sprintf("failed to load users: %v", err))

				return err
			}

			for _, user := range users {
				err = appMonitorService.NotifyUserAboutCertExpiration(user)

				if err != nil {
					logger.Error(fmt.Sprintf("failed to notify user %d about ssl expiration: %v", user.ID, err))
				}
			}

			return nil
		},
	}

	return notificationsCmd
}
