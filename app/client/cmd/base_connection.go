package cmd

import (
	"fmt"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/flash1nho/GophKeeper/certs"
	"github.com/flash1nho/GophKeeper/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func BaseConnection() (*grpc.ClientConn, error) {
	settings := config.Settings()
	certs, err := certs.NewCerts("client")

	if err != nil {
		return nil, fmt.Errorf("сертификаты не найдены: %v", err)
	}

	clientCert, err := tls.LoadX509KeyPair(certs.Cert, certs.Key)

	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить сертификат клиента: %v", err)
	}

	certCA, err := ioutil.ReadFile(certs.CertCA)

	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать CA сертификат: %v", err)
	}

	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(certCA) {
		return nil, fmt.Errorf("не удалось добавить CA сертификат: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	creds := credentials.NewTLS(tlsConfig)
	conn, err := grpc.Dial(settings.GrpcServerAddress, grpc.WithTransportCredentials(creds))

	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться: %v", err)
	}

	return conn, err
}
