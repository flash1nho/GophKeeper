package grpc

import (
	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/facade"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GrpcHandler struct {
	GrpcPublicHandler  *GrpcPublicHandler
	GrpcPrivateHandler *GrpcPrivateHandler
	Pool               *pgxpool.Pool
	Settings           config.SettingsObject
}

func NewGrpcHandler(pool *pgxpool.Pool, settings config.SettingsObject, facade *facade.Facade) *GrpcHandler {
	return &GrpcHandler{
		GrpcPublicHandler:  &GrpcPublicHandler{Pool: pool, Settings: settings},
		GrpcPrivateHandler: &GrpcPrivateHandler{Pool: pool, Settings: settings, facade: facade},
		Pool:               pool,
		Settings:           settings,
	}
}
