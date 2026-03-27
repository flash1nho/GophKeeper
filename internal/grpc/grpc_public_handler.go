package grpc

import (
	"context"
	"fmt"

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

	if req.Login == "" || req.Password == "" {
		return nil, fmt.Errorf("введите логин и пароль")
	}

	userID, err := users.Register(ctx, req.Login, req.Password, g.Pool)

	if err != nil {
		return nil, err
	}

	if userID == 0 {
		return nil, fmt.Errorf("логин уже занят")
	}

	token, err := CreateTokenFor(userID)

	if err != nil {
		return nil, err
	}

	response.Token = token

	return &response, nil
}

func (g *GrpcPublicHandler) Login(ctx context.Context, req *UserLoginRequest) (*UserLoginResponse, error) {
	var response UserLoginResponse

	if req.Login == "" || req.Password == "" {
		return nil, fmt.Errorf("введите логин и пароль")
	}

	userID, err := users.Login(ctx, req.Login, req.Password, g.Pool)

	if err != nil {
		return nil, err
	}

	if userID == 0 {
		return nil, fmt.Errorf("неверная логин или пароль")
	}

	token, err := CreateTokenFor(userID)

	if err != nil {
		return nil, err
	}

	response.Token = token

	return &response, nil
}
