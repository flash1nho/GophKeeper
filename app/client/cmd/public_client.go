package cmd

import (
	"fmt"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/flash1nho/GophKeeper/certs"
	"github.com/flash1nho/GophKeeper/internal/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func PublicClient() (pb.GophKeeperPublicServiceClient, *grpc.ClientConn, error) {
	certs, err := certs.NewCerts("client")

	if err != nil {
		return nil, nil, fmt.Errorf("сертификаты не найдены: %v", err)
	}

	clientCert, err := tls.LoadX509KeyPair(certs.Cert, certs.Key)

	if err != nil {
		return nil, nil, fmt.Errorf("не удалось загрузить сертификат клиента: %v", err)
	}

	certCA, err := ioutil.ReadFile(certs.CertCA)

	if err != nil {
		return nil, nil, fmt.Errorf("не удалось прочитать CA сертификат: %v", err)
	}

	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(certCA) {
		return nil, nil, fmt.Errorf("не удалось добавить CA сертификат: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	creds := credentials.NewTLS(tlsConfig)
	conn, err := grpc.Dial(config.GrpcServerAddress, grpc.WithTransportCredentials(creds))

	if err != nil {
		return nil, nil, fmt.Errorf("не удалось подключиться: %v", err)
	}

	client := pb.NewGophKeeperPublicServiceClient(conn)

	return client, conn, nil
}
