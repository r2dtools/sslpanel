package certificate

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/r2dtools/agentintegration"
	"golang.org/x/net/idna"
)

type SelfSignedCertificateReuest struct {
	CertName,
	CommonName,
	Email,
	Country,
	Province,
	Locality,
	Organization string
	AltNames []string
}

// GetX509CertificateFromRequest retrieves certificate from http request to domain
func GetX509CertificateFromRequest(domain string) ([]*x509.Certificate, error) {
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: time.Minute}, "tcp", domain+":443", &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		return nil, err
	}

	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
}

// ConvertX509CertificateToIntCert converts x509 certificate to agentintegration.Certificate
func ConvertX509CertificateToIntCert(certificate *x509.Certificate, roots []*x509.Certificate) *agentintegration.Certificate {
	certPool := x509.NewCertPool()

	for _, root := range roots {
		certPool.AddCert(root)
	}

	opts := x509.VerifyOptions{
		Roots: certPool,
	}
	_, err := certificate.Verify(opts)
	isValid := err == nil

	var DNSNames []string
	p := idna.New()
	CN, err := p.ToUnicode(certificate.Subject.CommonName)

	if err != nil {
		CN = certificate.Subject.CommonName
	}

	for _, name := range certificate.DNSNames {
		uName, err := p.ToUnicode(name)
		if err == nil {
			DNSNames = append(DNSNames, uName)
		}
	}

	cert := agentintegration.Certificate{
		DNSNames:       DNSNames,
		CN:             CN,
		EmailAddresses: certificate.EmailAddresses,
		Organization:   certificate.Subject.Organization,
		Country:        certificate.Subject.Country,
		Province:       certificate.Subject.Province,
		Locality:       certificate.Subject.Locality,
		IsCA:           certificate.IsCA,
		ValidFrom:      certificate.NotBefore.Format(time.RFC822Z),
		ValidTo:        certificate.NotAfter.Format(time.RFC822Z),
		Issuer: agentintegration.Issuer{
			CN:           certificate.Issuer.CommonName,
			Organization: certificate.Issuer.Organization,
		},
		IsValid: isValid,
	}

	return &cert
}

// GetCertificateForDomainFromRequest returns a certificate for a domain
func GetCertificateForDomainFromRequest(domain string) (*agentintegration.Certificate, error) {
	certs, err := GetX509CertificateFromRequest(domain)

	if err != nil {
		return nil, err
	}

	if len(certs) == 0 {
		return nil, nil
	}

	var roots []*x509.Certificate

	if len(certs) > 1 {
		roots = certs[1:]
	}

	return ConvertX509CertificateToIntCert(certs[0], roots), nil
}

// CreateSelfSignedCertificate creates self-signed certificate for the request
func CreateSelfSignedCertificate(certRequest *SelfSignedCertificateReuest) (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", err
	}

	subject := pkix.Name{
		CommonName: certRequest.CommonName,
	}
	if certRequest.Organization != "" {
		subject.Organization = []string{certRequest.Organization}
	}
	if certRequest.Country != "" {
		subject.Country = []string{certRequest.Country}
	}
	if certRequest.Locality != "" {
		subject.Locality = []string{certRequest.Locality}
	}
	if certRequest.Province != "" {
		subject.Province = []string{certRequest.Province}
	}
	template := x509.Certificate{
		SerialNumber:          big.NewInt(2022),
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              certRequest.AltNames,
	}

	if certRequest.Email != "" {
		template.EmailAddresses = []string{certRequest.Email}
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", err
	}
	certPem := new(bytes.Buffer)
	err = pem.Encode(certPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return "", err

	}
	privKeyPem := new(bytes.Buffer)
	err = pem.Encode(privKeyPem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return "", err
	}
	cert := strings.Join([]string{privKeyPem.String(), certPem.String()}, "\n")

	return cert, nil
}
