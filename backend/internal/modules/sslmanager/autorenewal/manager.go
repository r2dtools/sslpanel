package autorenewal

import (
	"backend/config"
	"backend/internal/app/panel/domain/dto"
	"backend/internal/app/panel/domain/provider"
	domainProvider "backend/internal/app/panel/domain/provider"
	domainStorage "backend/internal/app/panel/domain/storage"
	serverStorage "backend/internal/app/panel/server/storage"
	"backend/internal/modules/sslmanager/agent"
	"backend/internal/pkg/acme"
	serverAgent "backend/internal/pkg/agent"
	"backend/internal/pkg/logger"
	"fmt"

	"github.com/r2dtools/agentintegration"
)

type RenewLogWriter interface {
	WriteLog(serverID uint, successDomains []string, failedDomains map[string]error) error
}

const (
	defaultWorkersCount = 10
)

type RenewResult struct {
	ServerID       uint
	ServerName     string
	SuccessDomains []string
	FailedDomains  map[string]error
	Err            error
}

type AutoRenewalManager struct {
	config               *config.Config
	serverStorage        serverStorage.ServerStorage
	domainSettingStorage domainStorage.DomainSettingStorage
	domainProvider       domainProvider.DomainProvider
	logger               logger.Logger
	renewLogWriter       RenewLogWriter
}

type BlockReleaser interface {
	Release()
}

func (a AutoRenewalManager) Run(releaser <-chan struct{}) {
	defer func() {
		<-releaser
	}()

	servers, err := a.serverStorage.FindAll()

	if err != nil {
		a.logger.Error(fmt.Sprintf("renewal failed: %v", err))

		return
	}

	serversCount := len(servers)
	a.logger.Debug(fmt.Sprintf("found %d servers", serversCount))

	if serversCount == 0 {
		return
	}

	jobs := make(chan serverStorage.Server, serversCount)
	results := make(chan RenewResult, serversCount)

	for _, server := range servers {
		jobs <- server
	}

	close(jobs)

	workersCount := min(serversCount, defaultWorkersCount)

	for range workersCount {
		go a.renewWorker(jobs, results)
	}

	for range serversCount {
		result := <-results

		if result.Err != nil {
			a.logger.Error(fmt.Sprintf("renewal failed, server: %s, err: %v", result.ServerName, result.Err))

			continue
		}

		err = a.renewLogWriter.WriteLog(result.ServerID, result.SuccessDomains, result.FailedDomains)

		if err != nil {
			a.logger.Error(fmt.Sprintf("failed to write renew log: %v", err))
		}
	}

	close(results)
}

func (a AutoRenewalManager) renewWorker(
	servers <-chan serverStorage.Server,
	results chan<- RenewResult,
) {
	for server := range servers {
		domains, err := a.domainProvider.GetServerDomains(server.Guid)
		result := RenewResult{
			ServerID:   server.ID,
			ServerName: server.Name,
		}

		if err != nil {
			result.Err = err
			results <- result

			continue
		}

		a.logger.Debug(fmt.Sprintf("found %d domains on server %s", len(domains), server.Name))

		sAgent, err := serverAgent.NewAgent(
			server.Ipv4Address,
			server.Ipv6Address,
			server.Token,
			server.AgentPort,
			a.logger,
		)

		if err != nil {
			result.Err = err
			results <- result

			continue
		}

		certificateAgent := agent.NewCertificateAgent(sAgent)

		succeededDomains := []string{}
		failedDomains := map[string]error{}

		for _, domain := range domains {
			domainName := domain.ServerName
			cert := domain.Certificate

			if cert == nil {
				a.logger.Debug(fmt.Sprintf("skip renewal, no certificate, domain: %s", domainName))

				continue
			}

			var email string

			setting, err := a.domainSettingStorage.FindByDomain(domainName, server.Guid, "renewal")

			if err != nil {
				failedDomains[domainName] = err

				continue
			}

			if setting == nil || setting.SettingValue == "false" {
				a.logger.Info(fmt.Sprintf("auto renewal is disabled, server: %s, domain: %s", server.Name, domainName))

				continue
			}

			if len(cert.EmailAddresses) == 0 {
				emailSetting, err := a.domainSettingStorage.FindByDomain(domainName, server.Guid, "email")

				if err == nil && emailSetting != nil {
					email = emailSetting.SettingValue
				}
			} else {
				email = cert.EmailAddresses[0]
			}

			if !cert.IsLetsEncrypt() {
				a.logger.Debug(fmt.Sprintf("skip renewal for domain %s: not Let`s Encrypt certificate", domainName))

				continue
			}

			isAboutToExpire, err := cert.IsAboutToExpire(a.config.CertAboutToExpireInterval)

			if err != nil {
				failedDomains[domainName] = err

				continue
			}

			if !isAboutToExpire {
				a.logger.Debug(fmt.Sprintf("skip renewal for domain %s: certificate is actual", domainName))

				continue
			}

			err = issueCert(certificateAgent, email, domain)

			if err != nil {
				failedDomains[domainName] = err

				continue
			}

			succeededDomains = append(succeededDomains, domainName)
		}

		result.SuccessDomains = succeededDomains
		result.FailedDomains = failedDomains

		results <- result
	}
}

func issueCert(certificateAgent *agent.CertificateAgent, email string, domain dto.Domain) error {
	_, err := certificateAgent.Issue(agentintegration.CertificateIssueRequestData{
		Email:         email,
		ServerName:    domain.Certificate.CN,
		WebServer:     domain.WebServer,
		ChallengeType: acme.HttpChallengeType,
		Subjects:      domain.Certificate.DNSNames,
		Assign:        true,
	})

	return err
}

func CreateAutoRenewalManager(
	serverStorage serverStorage.ServerStorage,
	domainSettingStorage domainStorage.DomainSettingStorage,
	domainProvider provider.DomainProvider,
	config *config.Config,
	logger logger.Logger,
	renewLogWriter RenewLogWriter,
) AutoRenewalManager {
	return AutoRenewalManager{
		serverStorage:        serverStorage,
		domainProvider:       domainProvider,
		domainSettingStorage: domainSettingStorage,
		config:               config,
		logger:               logger,
		renewLogWriter:       renewLogWriter,
	}
}
