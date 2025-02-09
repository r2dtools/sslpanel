package service

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

type DomainSetting struct {
	ID           int    `json:"id"`
	SettingName  string `json:"setting_name"`
	SettingValue string `json:"setting_value"`
}

type DomainConfigRequest struct {
	ServerGuid string
	DomainName string
	WebServer  string `form:"webserver"`
}

type DomainRequest struct {
	ServerGuid string
	DomainName string
	WebServer  string `form:"webserver"`
}

type DomainCertificateRequest struct {
	ServerGuid string
	DomainName string
	WebServer  string `form:"webserver"`
}
