package config

import (
	"os"

	"embed"

	"github.com/flash1nho/GophKeeper/internal/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//go:embed all:.env*
var envFiles embed.FS

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

	data, err := envFiles.ReadFile(".env")

	if err != nil {
		data, _ = envFiles.ReadFile(".env.example")
	}

	env, err := godotenv.Unmarshal(string(data))

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
