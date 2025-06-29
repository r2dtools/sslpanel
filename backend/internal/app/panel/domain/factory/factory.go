package factory

import (
	"backend/internal/app/panel/domain/dto"
	"strconv"
	"strings"

	"github.com/r2dtools/agentintegration"
)

func CreateDomain(vhost agentintegration.VirtualHost) *dto.Domain {
	serverName := strings.Trim(vhost.ServerName, ".")
	serverNameParts := strings.Split(serverName, ".")

	// skip vhost names like 'domain'
	if len(serverNameParts) <= 1 {
		return nil
	}

	var addresses []dto.DomainAddress

	for _, address := range vhost.Addresses {
		port, err := strconv.Atoi(address.Port)

		if err != nil {
			continue
		}

		addresses = append(addresses, dto.DomainAddress{
			IsIpv6: address.IsIpv6,
			Host:   address.Host,
			Port:   port,
		})
	}

	return &dto.Domain{
		FilePath:    vhost.FilePath,
		ServerName:  vhost.ServerName,
		DocRoot:     vhost.DocRoot,
		WebServer:   vhost.WebServer,
		Aliases:     vhost.Aliases,
		Ssl:         vhost.Ssl,
		Addresses:   addresses,
		Certificate: CreateCertificate(vhost.Certificate),
	}
}

func CreateCertificate(cert *agentintegration.Certificate) *dto.DomainCertificate {
	if cert == nil {
		return nil
	}

	return &dto.DomainCertificate{
		CN:             cert.CN,
		ValidFrom:      cert.ValidFrom,
		ValidTo:        cert.ValidTo,
		DNSNames:       cert.DNSNames,
		EmailAddresses: cert.EmailAddresses,
		Organization:   cert.Organization,
		Country:        cert.Country,
		Locality:       cert.Locality,
		Province:       cert.Province,
		IsValid:        cert.IsValid,
		IsCA:           cert.IsCA,
		Issuer:         dto.Issuer(cert.Issuer),
	}
}
