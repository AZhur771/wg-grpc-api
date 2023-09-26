package certs

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"google.golang.org/grpc/credentials"
)

func getTLSConfig(host, caCertPath, certPath, keyPath string, certOpt tls.ClientAuthType) (*tls.Config, error) {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool
	if certOpt > tls.RequestClientCert {
		caCert, err = ioutil.ReadFile(caCertPath)
		if err != nil {
			log.Fatal("Error opening cert file", caCertPath, ", error ", err)
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ServerName: host,
		ClientAuth: certOpt,
		ClientCAs:  caCertPool,
		Certificates: []tls.Certificate{
			cert,
		},
		MinVersion: tls.VersionTLS12,
	}, nil
}

func LoadTLSCredentials(host, caCertPath, certPath, keyPath string, certOpt tls.ClientAuthType) (credentials.TransportCredentials, error) {
	config, err := getTLSConfig(host, caCertPath, certPath, keyPath, certOpt)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(config), nil
}
