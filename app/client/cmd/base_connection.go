package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"time"

	"github.com/flash1nho/GophKeeper/certs"
	"github.com/flash1nho/GophKeeper/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	ErrCertsNotFound    = errors.New("сертификаты не найдены")
	ErrLoadClientCert   = errors.New("не удалось загрузить сертификат клиента")
	ErrReadCACert       = errors.New("не удалось прочитать CA сертификат")
	ErrAppendCACert     = errors.New("не удалось добавить CA сертификат")
	ErrConnectionFailed = errors.New("не удалось подключиться")
)

type jwtCredentials struct {
	token string
}

func (j jwtCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

func (j jwtCredentials) RequireTransportSecurity() bool {
	return true
}

func BaseConnection(settings config.SettingsObject, token string) (*grpc.ClientConn, error) {
	certs, err := certs.NewCerts("client")
	if err != nil {
		return nil, ErrCertsNotFound
	}

	clientCert, err := tls.LoadX509KeyPair(certs.Cert, certs.Key)
	if err != nil {
		return nil, ErrLoadClientCert
	}

	certCA, err := os.ReadFile(certs.CertCA)
	if err != nil {
		return nil, ErrReadCACert
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(certCA) {
		return nil, ErrAppendCACert
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	creds := credentials.NewTLS(tlsConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	}

	if token != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(jwtCredentials{token: token}))
	}

	conn, err := grpc.DialContext(ctx, settings.GrpcServerAddress, opts...)

	if err != nil {
		return nil, ErrConnectionFailed
	}

	return conn, nil
}
