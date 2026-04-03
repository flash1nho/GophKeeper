package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/flash1nho/GophKeeper/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/flash1nho/GophKeeper/internal/certs"
)

var (
	ErrLoadClientCert   = errors.New("не удалось загрузить сертификат клиента")
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
	clientCert, err := tls.X509KeyPair(certs.ClientCrt, certs.ClientKey)

	if err != nil {
		return nil, ErrLoadClientCert
	}

	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(certs.CaCrt) {
		return nil, ErrAppendCACert
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	creds := credentials.NewTLS(tlsConfig)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	if token != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(jwtCredentials{token: token}))
	}

	conn, err := grpc.NewClient(settings.GrpcServerAddress, opts...)

	if err != nil {
		return nil, ErrConnectionFailed
	}

	return conn, nil
}
