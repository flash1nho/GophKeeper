package grpc

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type GrpcHandler struct {
	GrpcPublicHandler  *GrpcPublicHandler
	GrpcPrivateHandler *GrpcPrivateHandler
	Pool               *pgxpool.Pool
	Log                *zap.Logger
}

func NewGrpcHandler(pool *pgxpool.Pool, log *zap.Logger) *GrpcHandler {
	return &GrpcHandler{
		GrpcPublicHandler:  &GrpcPublicHandler{Pool: pool, Log: log},
		GrpcPrivateHandler: &GrpcPrivateHandler{Pool: pool, Log: log},
		Pool:               pool,
		Log:                log,
	}
}
