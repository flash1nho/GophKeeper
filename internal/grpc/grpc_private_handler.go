package grpc

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type GrpcPrivateHandler struct {
	UnimplementedGophKeeperPrivateServiceServer

	Pool *pgxpool.Pool
	Log  *zap.Logger
}
