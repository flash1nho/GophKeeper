package config

import (
	"github.com/flash1nho/GophKeeper/internal/logger"

	"go.uber.org/zap"
)

const (
	DatabaseURI       = "postgres://gophkeeper:gophkeeper@localhost:5433/gophkeeper?sslmode=disable"
	CryptoKey         = "NzFKH^>h*a{pkCom"
	GrpcServerAddress = "localhost:3200"
)

type SettingsObject struct {
	DatabaseURI       string
	CryptoKey         string
	GrpcServerAddress string
	Log               *zap.Logger
}

func Settings() SettingsObject {
	logger.Initialize("info")

	return SettingsObject{
		DatabaseURI:       DatabaseURI,
		CryptoKey:         CryptoKey,
		GrpcServerAddress: GrpcServerAddress,
		Log:               logger.Log,
	}
}
