package config

import (
	"os"

	_ "embed"

	"github.com/flash1nho/GophKeeper/internal/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//go:embed .env
var envFile string

type SettingsObject struct {
	DatabaseDSN       string
	MasterKey         []byte
	GrpcServerAddress string
	Log               *zap.Logger
}

func Settings() SettingsObject {
	err := logger.Initialize("info")

	if err != nil {
		logger.Log.Fatal("Ошибка загрузки logger")
	}

	env, err := godotenv.Unmarshal(envFile)

	if err != nil {
		logger.Log.Fatal("Ошибка парсинга встроенного .env файла", zap.Error(err))
	}

	for key, value := range env {
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}

	return SettingsObject{
		DatabaseDSN:       os.Getenv("DATABASE_DSN"),
		GrpcServerAddress: os.Getenv("GRPC_SERVER_ADDRESS"),
		MasterKey:         []byte(os.Getenv("MASTER_KEY")),
		Log:               logger.Log,
	}
}
