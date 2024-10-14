package cmd

import (
	"backend/certificate"
	"backend/logger"
	"backend/modules/certificatemonitor/models"
	"fmt"

	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh certificate data for sites",
	RunE: func(cmd *cobra.Command, args []string) error {
		sites, err := models.GetAllSites()

		if err != nil {
			logger.Debug(fmt.Sprintf("could not load sites '%v' ...", err))
			return err
		}

		for _, site := range sites {
			logger.Debug(fmt.Sprintf("refreshing site data '%s' ...", site.URL))
			cert, err := certificate.GetCertificateForDomainFromRequest(site.URL)

			if err != nil {
				logger.Error(fmt.Sprintf("could not get site '%s' certificate: %v", site.URL, err))
				continue
			}

			if err = site.UpdateCertData(cert); err != nil {
				logger.Error(fmt.Sprintf("could not update site '%s' certificate: %v", site.URL, err))
				continue
			}
			logger.Debug(fmt.Sprintf("site data '%s' successfully refreshed", site.URL))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
