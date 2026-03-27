package config

import (
	"os"

	"github.com/flash1nho/GophKeeper/internal/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type SettingsObject struct {
	DatabaseDSN       string
	CryptoKey         string
	GrpcServerAddress string
	Log               *zap.Logger
}

func Settings() SettingsObject {
	logger.Initialize("info")

	err := godotenv.Load()

	if err != nil {
		logger.Log.Fatal("Ошибка загрузки .env файла")
	}

	return SettingsObject{
		DatabaseDSN:       os.Getenv("DATABASE_DSN"),
		GrpcServerAddress: os.Getenv("GRPC_SERVER_ADDRESS"),
		CryptoKey:         os.Getenv("CRYPTO_KEY"),
		Log:               logger.Log,
	}
}
