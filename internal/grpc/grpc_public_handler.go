package grpc

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flash1nho/GophKeeper/internal/models/users"
	"go.uber.org/zap"
)

type GrpcPublicHandler struct {
	UnimplementedGophKeeperPublicServiceServer

	Pool *pgxpool.Pool
	Log  *zap.Logger
}

func (g *GrpcPublicHandler) Register(ctx context.Context, req *UserRegisterRequest) (*UserRegisterResponse, error) {
	var response UserRegisterResponse

	user := users.NewUser(req.Login, req.Password)
	token, err := user.UserRegister(ctx, g.Pool)

	if err != nil {
		return nil, err
	}

	response.Token = token

	return &response, nil
}

func (g *GrpcPublicHandler) Login(ctx context.Context, req *UserLoginRequest) (*UserLoginResponse, error) {
	var response UserLoginResponse

	user := users.NewUser(req.Login, req.Password)
	token, err := user.UserLogin(ctx, g.Pool)

	if err != nil {
		return nil, err
	}

	response.Token = token

	return &response, nil
}
