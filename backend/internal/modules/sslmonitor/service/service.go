package service

import (
	userStorage "backend/internal/app/panel/user/storage"
	domainStorage "backend/internal/modules/sslmonitor/storage"
	"backend/internal/pkg/certificate"
	"backend/internal/pkg/logger"
	"backend/internal/pkg/notification"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"time"
)

const aboutToExpireDaysEdge = 30
const notificationSubject = "Certificate expiration notice"

var ErrDomainNotFound = errors.New("domain not found")

type MonitorService struct {
	domainStorage       domainStorage.DomainStorage
	notificationService notification.EmailNotificationService
	logger              logger.Logger
}

func (s MonitorService) FindUserDomains(userID int) ([]*Domain, error) {
	domains, err := s.domainStorage.FindByUserID(userID)

	if err != nil {
		return nil, err
	}

	var appDomains []*Domain

	for _, domain := range domains {
		appDomains = append(appDomains, createAppDomain(domain))
	}

	return appDomains, nil
}

func (s MonitorService) RemoveDomain(id int) error {
	domain, err := s.domainStorage.FindByID(id)

	if err != nil {
		return err
	}

	if domain == nil {
		return ErrDomainNotFound
	}

	return s.domainStorage.Remove(domain)
}

func (s MonitorService) AddDomain(requestData AddDomainRequest, userID int) error {
	domainModel := domainStorage.Domain{
		URL:    requestData.URL,
		Data:   requestData.Data,
		UserID: uint(userID),
	}

	return s.domainStorage.Save(&domainModel)
}

func (s MonitorService) RefreshDomain(id int) (*Domain, error) {
	domain, err := s.domainStorage.FindByID(id)

	if err != nil {
		return nil, err
	}

	if domain == nil {
		return nil, ErrDomainNotFound
	}

	cert, err := certificate.GetCertificateForDomainFromRequest(domain.URL)

	if err != nil {
		return nil, err
	}

	err = s.domainStorage.UpdateCertData(domain, cert)

	if err != nil {
		return nil, err
	}

	return createAppDomain(domain), nil
}

func (s MonitorService) NotifyUserAboutCertExpiration(user *userStorage.User) error {
	domains, err := s.domainStorage.FindByUserID(int(user.ID))

	if err != nil {
		return err
	}

	if len(domains) == 0 {
		return nil
	}

	certExpireDays := make(map[string]int)

	type tplData struct {
		Expired, AboutToExpire map[string]int
		AboutToExpireDaysEdge  int
	}

	for _, domain := range domains {
		if err := s.collectDomainCertExpirationDays(domain, certExpireDays); err != nil {
			s.logger.Error("failed to collect domain %s cert expiration days: %v", domain.URL, err)

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
		s.logger.Debug("all certificates are up to date. skip notifications sendindg.")

		return nil
	}

	data := tplData{expired, aboutToExpire, aboutToExpireDaysEdge}
	tplPath := filepath.Join("certificatemonitor", "notification-template")

	return s.notificationService.CreateAndSendPlainNotification("certNotificaton", tplPath, user.Email, notificationSubject, data)
}

func (s MonitorService) collectDomainCertExpirationDays(domain *domainStorage.Domain, certExpireDays map[string]int) error {
	data, err := domain.GetData()

	if err != nil {
		return fmt.Errorf("failed to get domain '%s' data: %v", domain.URL, err)
	}

	cert := data.Cert

	if cert == nil {
		s.logger.Debug(fmt.Sprintf("domain '%s' does not have a certificate. Skip check.", domain.URL))

		return nil
	}

	validToTime, err := time.Parse(time.RFC822Z, cert.ValidTo)

	if err != nil {
		return fmt.Errorf("failed to parse certificate expiration date '%s' data: %v", cert.ValidTo, err)
	}

	expireDays := math.Floor(validToTime.Sub(time.Now()).Hours() / 24)
	certExpireDays[domain.URL] = int(expireDays)

	return nil
}

func createAppDomain(domain *domainStorage.Domain) *Domain {
	return &Domain{
		ID:        int(domain.ID),
		URL:       domain.URL,
		Data:      domain.Data,
		UserID:    int(domain.UserID),
		CreatedAt: domain.CreatedAt,
	}
}

func NewMonitorService(domainStorage domainStorage.DomainStorage) MonitorService {
	return MonitorService{
		domainStorage: domainStorage,
	}
}
