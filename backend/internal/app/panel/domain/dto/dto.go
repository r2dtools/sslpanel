package dto

import (
	"strings"
	"time"
)

type Domain struct {
	FilePath    string             `json:"filepath"`
	ServerName  string             `json:"servername"`
	DocRoot     string             `json:"docroot"`
	WebServer   string             `json:"webserver"`
	Aliases     []string           `json:"aliases"`
	Ssl         bool               `json:"ssl"`
	Addresses   []DomainAddress    `json:"addresses"`
	Certificate *DomainCertificate `json:"certificate"`
}

type DomainAddress struct {
	IsIpv6 bool   `json:"isIpv6"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

type DomainCertificate struct {
	CN             string   `json:"cn"`
	ValidFrom      string   `json:"validfrom"`
	ValidTo        string   `json:"validto"`
	DNSNames       []string `json:"dnsnames"`
	EmailAddresses []string `json:"emailaddresses"`
	Organization   []string `json:"organization"`
	Province       []string `json:"province"`
	Country        []string `json:"country"`
	Locality       []string `json:"locality"`
	IsCA           bool     `json:"isca"`
	IsValid        bool     `json:"isvalid"`
	Issuer         Issuer   `json:"issuer"`
}

func (c DomainCertificate) IsLetsEncrypt() bool {
	for _, org := range c.Organization {
		if strings.Contains(org, "Let's Encrypt") || strings.Contains(org, "good guys") {
			return true
		}
	}

	for _, org := range c.Issuer.Organization {
		if strings.Contains(org, "Let's Encrypt") || strings.Contains(org, "good guys") {
			return true
		}
	}

	return strings.Contains(c.CN, "Let's Encrypt") || strings.Contains(c.CN, "R3")
}

func (c DomainCertificate) IsSelfSigned() bool {
	return c.CN == c.Issuer.CN
}

func (c DomainCertificate) IsAboutToExpire(duration time.Duration) (bool, error) {
	validTo, err := time.Parse(time.RFC822Z, c.ValidTo)

	if err != nil {
		return false, err
	}

	return time.Until(validTo) < duration, nil
}

type Issuer struct {
	CN           string   `json:"cn"`
	Organization []string `json:"organization"`
}
