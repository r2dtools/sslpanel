package service

import "time"

type Server struct {
	ID           int       `json:"id"`
	Guid         string    `json:"guid"`
	Name         string    `json:"name"`
	OsCode       string    `json:"os_code"`
	OsVersion    string    `json:"os_version"`
	Ipv4Address  string    `json:"ipv4_address"`
	Ipv6Address  string    `json:"ipv6_address"`
	AgentVersion string    `json:"agent_version"`
	AgentPort    int       `json:"agent_port"`
	IsActive     int       `json:"is_active"`
	IsRegistered int       `json:"is_registered"`
	AccountID    int       `json:"account_id"`
	CreatedAt    time.Time `json:"created_at"`
	Token        string    `json:"token"`
}

type ServerDetails struct {
	Server

	HostName       string   `json:"hostname"`
	Os             string   `json:"os"`
	PlatformFamily string   `json:"platform_family"`
	KernelVersion  string   `json:"kernal_version"`
	KernelArch     string   `json:"kernal_arch"`
	Virtualization string   `json:"virtualization"`
	Uptime         uint64   `json:"uptime"`
	BootTime       uint64   `json:"boottime"`
	Domains        []Domain `json:"domains"`
}

type NewServerRequest struct {
	Name        string `json:"name" validate:nonzero`
	Ipv4Address string `json:"ipv4_address"`
	Ipv6Address string `json:"ipv6_address"`
	AgentPort   int    `json:"agent_port" validate:nonzero`
	Token       string `json:"token" validate:nonzero`
	AccountID   int
}

type UpdateServerRequest struct {
	ID          int    `json:"id" validate:nonzero`
	Name        string `json:"name" validate:nonzero`
	Ipv4Address string `json:"ipv4_address"`
	Ipv6Address string `json:"ipv6_address"`
	AgentPort   int    `json:"agent_port" validate:nonzero`
	Token       string `json:"token" validate:nonzero`
}

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

type Issuer struct {
	CN           string   `json:"cn"`
	Organization []string `json:"organization"`
}
