package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/flash1nho/GophKeeper/internal/authenticator"
)

type grpcProvider struct{}

func (p *grpcProvider) GetCookie(ctx context.Context, _ string) (string, error) {
	var cookie string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("authorization")

		if len(values) > 0 {
			cookie = values[0]
		}
	}

	return cookie, nil
}

func (p *grpcProvider) SetCookie(ctx context.Context, cookieName, cookieValue string) error {
	header := metadata.Pairs("set-cookie", cookieName+"="+cookieValue+"; Path=/")
	return grpc.SendHeader(ctx, header)
}

func Auth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	auth := authenticator.NewAuthenticator()
	ctx, err := auth.Authenticate(ctx, &grpcProvider{})

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return handler(ctx, req)
}

func CreateSignedCookie(ctx context.Context, userID string) error {
	auth := authenticator.NewAuthenticator()
	err := auth.CreateSignedCookie(ctx, userID, &grpcProvider{})

	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
