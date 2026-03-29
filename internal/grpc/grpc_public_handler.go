package grpc

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/models/users"
)

type GrpcPublicHandler struct {
	UnimplementedGophKeeperPublicServiceServer

	Pool     *pgxpool.Pool
	Settings config.SettingsObject
}

func (g *GrpcPublicHandler) Register(ctx context.Context, req *UserRegisterRequest) (*UserRegisterResponse, error) {
	user := users.NewUser(req.Login, req.Password)
	token, err := user.UserRegister(ctx, g.Pool, g.Settings, req.Secret)

	if err != nil {
		return nil, err
	}

	return &UserRegisterResponse{Token: token}, nil
}

func (g *GrpcPublicHandler) Login(ctx context.Context, req *UserLoginRequest) (*UserLoginResponse, error) {
	user := users.NewUser(req.Login, req.Password)
	token, err := user.UserLogin(ctx, g.Pool, g.Settings)

	if err != nil {
		return nil, err
	}

	return &UserLoginResponse{Token: token}, nil
}
