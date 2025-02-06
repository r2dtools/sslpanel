package service

type SelfSignedCertificateRequest struct {
	CertName,
	CommonName,
	Email,
	Country,
	Province,
	Locality,
	Organization string
	AltNames []string
}

type CertificateIssueRequest struct {
	Email            string            `json:"email"`
	ServerName       string            `json:"servername"`
	DocRoot          string            `json:"docroot"`
	WebServer        string            `json:"webserver"`
	ChallengeType    string            `json:"challengetype"`
	Subjects         []string          `json:"subjects"`
	AdditionalParams map[string]string `json:"params"`
	Assign           bool              `json:"assign"`
}

type CommonDirStatusRequest struct {
	ServerName string `json:"servername"`
	WebServer  string `json:"webserver"`
}

type CommonDirStatusChangeRequest struct {
	Status     bool   `json:"status"`
	ServerName string `json:"servername"`
	WebServer  string `json:"webserver"`
}

type CommonDirStatusResponse struct {
	Status bool `json:"status"`
}
