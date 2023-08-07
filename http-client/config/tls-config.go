package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"strings"
)

const certPath = "./certs/"

func TransportConfig(insecure bool) (*http.Transport, error) {
	rootCAs, err := loadCertPool()

	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		InsecureSkipVerify: insecure,
		RootCAs:            rootCAs,
		Certificates:       []tls.Certificate{loadClientCert()},
	}

	return &http.Transport{
		TLSClientConfig: config,
	}, nil
}

func loadCertPool() (*x509.CertPool, error) {

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	for _, cert := range loadLocalCerts() {
		file, _ := os.ReadFile(cert)

		if ok := rootCAs.AppendCertsFromPEM(file); !ok {
			return nil, errors.New("Failed to append cert " + cert)
		}
	}

	return rootCAs, nil
}

func loadLocalCerts() []string {
	certDir, err := os.ReadDir(certPath)

	if err != nil {
		panic("Certs path not found")
	}

	allCerts := []string{}

	for _, certFile := range certDir {

		if !strings.HasSuffix(certFile.Name(), ".pem") {
			continue
		}

		allCerts = append(allCerts, certPath+certFile.Name())
	}

	return allCerts
}

func loadClientCert() tls.Certificate {
	certificate, _ := tls.LoadX509KeyPair(
		certPath+"client.crt",
		certPath+"client.key",
	)

	return certificate
}
