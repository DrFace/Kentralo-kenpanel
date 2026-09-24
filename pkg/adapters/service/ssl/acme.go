package ssl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

// ACMEManager manages automated TLS certificate requests and HTTP-01 challenge files.
type ACMEManager struct {
	webRootDir string
	certDir    string
}

func NewACMEManager(webRootDir, certDir string) *ACMEManager {
	return &ACMEManager{
		webRootDir: webRootDir,
		certDir:    certDir,
	}
}

// WriteHTTP01Challenge writes the ACME challenge token to .well-known/acme-challenge/
func (m *ACMEManager) WriteHTTP01Challenge(domain, token, keyAuth string) (string, error) {
	challengeDir := filepath.Join(m.webRootDir, domain, ".well-known", "acme-challenge")
	if err := os.MkdirAll(challengeDir, 0755); err != nil {
		return "", fmt.Errorf("failed creating acme challenge directory: %w", err)
	}

	filePath := filepath.Join(challengeDir, token)
	if err := os.WriteFile(filePath, []byte(keyAuth), 0644); err != nil {
		return "", fmt.Errorf("failed writing acme challenge token: %w", err)
	}

	return filePath, nil
}

// CleanHTTP01Challenge removes the temporary challenge file after verification.
func (m *ACMEManager) CleanHTTP01Challenge(domain, token string) {
	filePath := filepath.Join(m.webRootDir, domain, ".well-known", "acme-challenge", token)
	_ = os.Remove(filePath)
}

// GenerateSelfSignedCertificate creates a fallback SSL certificate for initial bootstrap.
func (m *ACMEManager) GenerateSelfSignedCertificate(domain string) (string, string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"KenPanel Self-Signed"},
			CommonName:   domain,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{domain},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return "", "", err
	}

	if err := os.MkdirAll(m.certDir, 0700); err != nil {
		return "", "", err
	}

	certPath := filepath.Join(m.certDir, domain+".crt")
	keyPath := filepath.Join(m.certDir, domain+".key")

	certOut, err := os.Create(certPath)
	if err != nil {
		return "", "", err
	}
	defer certOut.Close()
	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", "", err
	}
	defer keyOut.Close()
	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certPath, keyPath, nil
}
