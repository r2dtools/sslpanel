package service

type Certificate struct {
	CN           string   `json:"cn"`
	ValidFrom    string   `json:"validfrom"`
	ValidTo      string   `json:"validto"`
	DNSNames     []string `json:"dnsnames"`
	Emails       []string `json:"emails"`
	Organization []string `json:"organization"`
	Province     []string `json:"province"`
	Country      []string `json:"country"`
	Locality     []string `json:"locality"`
	IsCA         bool     `json:"isca"`
	IsValid      bool     `json:"isvalid"`
	Issuer       Issuer   `json:"issuer"`
}

type Issuer struct {
	CN           string   `json:"cn"`
	Organization []string `json:"organization"`
}

type SelfSignedCertificateRequest struct {
	ServerGuid   string
	CertName     string   `json:"certName"`
	CommonName   string   `json:"commonName"`
	Email        string   `json:"email"`
	Country      string   `json:"country"`
	Province     string   `json:"province"`
	Locality     string   `json:"locality"`
	Organization string   `json:"organization"`
	AltNames     []string `json:"altNames"`
}

type CertificateIssueRequest struct {
	ServerGuid       string
	DomainName       string
	Email            string            `json:"email"`
	WebServer        string            `json:"webserver"`
	ChallengeType    string            `json:"challengetype"`
	Subjects         []string          `json:"subjects"`
	AdditionalParams map[string]string `json:"params"`
	Assign           bool              `json:"assign"`
}

type CommonDirStatusRequest struct {
	ServerGuid string
	DomainName string
	WebServer  string `form:"webserver"`
}

type ChangeCommonDirStatusRequest struct {
	ServerGuid string
	DomainName string
	Status     bool   `json:"status"`
	WebServer  string `json:"webserver"`
}

type CommonDirStatusResponse struct {
	Status bool `json:"status"`
}

type AssignCertificateRequest struct {
	ServerGuid string
	DomainName string
	WebServer  string `json:"webserver"`
	CertName   string `json:"name"`
}

type CertificatesRequest struct {
	Guid string
}

type CertificateUploadToStorageRequest struct {
	ServerGuid     string
	CertName       string
	PemCertificate string
}
