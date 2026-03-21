package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/flash1nho/GophKeeper/internal/authenticator"
)

type grpcProvider struct{}

func (p *grpcProvider) SetToken(ctx context.Context, userID int) error {
	// TODO
	return nil
}

func (p *grpcProvider) ParseToken(ctx context.Context, token string) error {
	// TODO
	return nil
}

func Auth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// TODO
	return true, nil
}

func CreateTokenFor(userID int) (string, error) {
	auth := authenticator.NewAuthenticator()
	token, err := auth.CreateToken(userID)

	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return token, nil
}
