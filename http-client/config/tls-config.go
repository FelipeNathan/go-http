package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func TransportConfig(insecure bool, certPath string) (*http.Transport, error) {
	rootCAs, err := loadCertPool(certPath)

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

func loadCertPool(certPath string) (*x509.CertPool, error) {

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	for _, cert := range loadLocalCerts(certPath) {
		file, _ := os.ReadFile(cert)

		if ok := rootCAs.AppendCertsFromPEM(file); !ok {
			return nil, errors.New("Failed to append cert " + cert)
		}
	}

	return rootCAs, nil
}

func loadLocalCerts(certPath string) []string {
	certDir, err := os.ReadDir(certPath)

	if err != nil {
		panic("Certs path not found")
	}

	allCerts := []string{}

	for _, certFile := range certDir {

		if !strings.HasSuffix(certFile.Name(), ".pem") {
			continue
		}

		fmt.Println(certFile.Name())

		allCerts = append(allCerts, certPath+certFile.Name())
	}

	return allCerts
}

func loadClientCert() tls.Certificate {
	certificate, _ := tls.LoadX509KeyPair(
		"./certs/client.crt",
		"./certs/client.key",
	)

	return certificate
}
