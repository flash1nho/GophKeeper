package grpc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flash1nho/GophKeeper/internal/authenticator"
	"github.com/flash1nho/GophKeeper/internal/models"
)

type GrpcHandler struct {
	UnimplementedGophKeeperServiceServer

	Pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *GrpcHandler {
	return &GrpcHandler{
		Pool: pool,
	}
}

func (g *GrpcHandler) Register(ctx context.Context, req *UserRegisterRequest) (*UserRegisterResponse, error) {
	var response UserRegisterResponse

	fmt.Println(req.Login)

	if req.Login == "" || req.Password == "" {
		err := fmt.Errorf("введите логин и пароль")
		return nil, err
	}

	userID, err := models.UserRegister(ctx, req.Login, req.Password, g.Pool)

	if err != nil {
		return nil, err
	}

	if userID == 0 {
		err = fmt.Errorf("логин уже занят")
		return nil, err
	}

	strUserID := strconv.Itoa(userID)
	err = CreateSignedCookie(ctx, strUserID)

	if err != nil {
		return nil, err
	}

	response.UserID = strUserID

	return &response, nil
}

func getUserFromContext(ctx context.Context) (string, error) {
	userID, err := authenticator.FromContext(ctx)

	if err != nil {
		return "", nil
	}

	return userID, nil
}
