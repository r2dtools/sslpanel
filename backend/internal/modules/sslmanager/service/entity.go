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
