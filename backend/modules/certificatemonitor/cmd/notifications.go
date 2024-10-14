package cmd

import (
	"backend/logger"
	mModels "backend/models"
	"backend/modules/certificatemonitor/models"
	"backend/notification"
	"fmt"
	"math"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const aboutToExpireDaysEdge = 30
const notificationSubject = "Certificate expiration notice"

var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Send notifications if nesseccary",
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := mModels.GetAllUsers()

		if err != nil {
			logger.Debug(fmt.Sprintf("could not load users: %v", err))
			return err
		}

		for _, user := range users {
			sites, err := models.GetUserSites(user.ID)

			if err != nil {
				logger.Error(fmt.Sprintf("could not load sites for the user with ID '%d': %v", user.ID, err))
				continue
			}

			if len(sites) == 0 {
				continue
			}

			if err = sendCertExpirationNotification(sites, user); err != nil {
				logger.Error(err.Error())
			}
		}

		return nil
	},
}

func sendCertExpirationNotification(sites []*models.Site, user *mModels.User) error {
	certExpireDays := make(map[string]int)

	type tplData struct {
		Expired, AboutToExpire map[string]int
		AboutToExpireDaysEdge  int
	}

	for _, site := range sites {
		if err := collectSiteCertExpirationDays(site, certExpireDays); err != nil {
			logger.Error(err.Error())
			continue
		}
	}

	expired := make(map[string]int)
	aboutToExpire := make(map[string]int)

	for url, days := range certExpireDays {
		if days < 0 {
			expired[url] = days
		} else if days < aboutToExpireDaysEdge {
			aboutToExpire[url] = days
		}
	}

	if len(expired) == 0 && len(aboutToExpire) == 0 {
		logger.Debug("all certificates are up to date. skip notifications sendindg.")
		return nil
	}

	data := tplData{expired, aboutToExpire, aboutToExpireDaysEdge}
	tplPath := filepath.Join("certificatemonitor", "notification-template")
	notification := &notification.EmailNotification{}

	if err := notification.CreateAndSendPlainNotification("certNotificaton", tplPath, user.Email, notificationSubject, data); err != nil {
		return err
	}

	return nil
}

func collectSiteCertExpirationDays(site *models.Site, certExpireDays map[string]int) error {
	sData, err := site.GetData()

	if err != nil {
		return fmt.Errorf("could not get site '%s' data: %v", site.URL, err)
	}

	cert := sData.Cert

	if cert == nil {
		logger.Debug(fmt.Sprintf("site '%s' does not have a certificate. Skip check.", site.URL))
		return nil
	}

	validToTime, err := time.Parse(time.RFC822Z, cert.ValidTo)

	if err != nil {
		return fmt.Errorf("could not parse certificate expiration date '%s' data: %v", cert.ValidTo, err)
	}

	expireDays := math.Floor(validToTime.Sub(time.Now()).Hours() / 24)
	certExpireDays[site.URL] = int(expireDays)

	return nil
}

func init() {
	rootCmd.AddCommand(notificationsCmd)
}
