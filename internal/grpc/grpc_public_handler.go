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
	var response UserRegisterResponse

	user := users.NewUser(0, req.Login, req.Password, req.Secret)
	token, err := user.UserRegister(ctx, g.Pool, g.Settings)

	if err != nil {
		return nil, err
	}

	response.Token = token

	return &response, nil
}

func (g *GrpcPublicHandler) Login(ctx context.Context, req *UserLoginRequest) (*UserLoginResponse, error) {
	var response UserLoginResponse

	user := users.NewUser(0, req.Login, req.Password, "")
	token, err := user.UserLogin(ctx, g.Pool, g.Settings)

	if err != nil {
		return nil, err
	}

	response.Token = token

	return &response, nil
}
