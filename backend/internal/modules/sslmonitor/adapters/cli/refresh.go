package cli

import (
	"backend/internal/modules/sslmonitor/service"
	"backend/internal/modules/sslmonitor/storage"
	"backend/internal/pkg/logger"
	"fmt"

	"github.com/spf13/cobra"
)

func getRefreshCmd(
	monitorService service.MonitorService,
	domainStorage storage.DomainStorage,
	logger logger.Logger,
) *cobra.Command {
	var refreshCmd = &cobra.Command{
		Use:   "refresh",
		Short: "Refresh certificate data for domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			domains, err := domainStorage.FindAll()

			if err != nil {
				logger.Debug(fmt.Sprintf("could not load domains '%v' ...", err))

				return err
			}

			for _, domain := range domains {
				logger.Debug(fmt.Sprintf("refreshing site data '%s' ...", domain.URL))

				_, err = monitorService.RefreshDomain(int(domain.ID))

				if err != nil {
					logger.Error("failed to refresh domain %s data: %v", domain.URL, err)

					continue
				}

				logger.Debug(fmt.Sprintf("domain data '%s' successfully refreshed", domain.URL))
			}

			return nil
		},
	}

	return refreshCmd
}
