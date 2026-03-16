package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"dario.cat/mergo"
	"github.com/flash1nho/GophKeeper/internal/logger"
	"go.uber.org/zap"
)

const (
	DefaultHost = "localhost:8080"
)

// Config — единая структура для всех источников
type Config struct {
	ServerAddress string `json:"server_address" env:"SERVER_ADDRESS"`
	DatabaseDSN   string `json:"database_dsn" env:"DATABASE_DSN"`
	EnableHTTPS   bool   `json:"enable_https" env:"ENABLE_HTTPS"`
	ConfigPath    string `json:"-" env:"CONFIG"`
}

type SettingsObject struct {
	Server      Server
	Log         *zap.Logger
	DatabaseDSN string
	EnableHTTPS bool
}

type Server struct {
	Addr string
}

func Settings() SettingsObject {
	logger.Initialize("info")

	// 1. Конфигурация из Флагов
	flagCfg := parseFlags()

	// 2. Конфигурация из ENV
	envCfg := parseEnv()

	// 3. Конфигурация из JSON (если путь указан)
	configPath := flagCfg.ConfigPath

	if envCfg.ConfigPath != "" {
		configPath = envCfg.ConfigPath
	}

	jsonCfg := parseJSON(configPath)

	// Итоговая сборка с помощью mergo.
	// Приоритет (от низшего к высшему): JSON -> ENV -> Flags
	finalCfg := jsonCfg

	// Накладываем ENV на JSON
	if err := mergo.Merge(&finalCfg, envCfg, mergo.WithOverride); err != nil {
		logger.Log.Error(fmt.Sprintf("Mergo error (ENV): %v", err))
	}

	// Накладываем Flags на результат (флаги имеют высший приоритет, если они установлены)
	// Для корректной работы mergo с флагами, в parseFlags нужно возвращать только заполненные значения
	if err := mergo.Merge(&finalCfg, flagCfg, mergo.WithOverride); err != nil {
		logger.Log.Error(fmt.Sprintf("Mergo error (Flags): %v", err))
	}

	// Дефолтные значения, если всё пусто
	if finalCfg.ServerAddress == "" {
		finalCfg.ServerAddress = DefaultHost
	}

	return SettingsObject{
		Server:        Server{Addr: finalCfg.ServerAddress},
		Log:           logger.Log,
		DatabaseDSN:   finalCfg.DatabaseDSN,
		EnableHTTPS:   finalCfg.EnableHTTPS,
		TrustedSubnet: finalCfg.TrustedSubnet,
	}
}

func parseFlags() Config {
	var c Config
	// Используем временные переменные, чтобы mergo не затер пустые строки дефолтами флагов
	serverAddress := flag.String("a", "", "значение может быть таким: "+DefaultHost)
	dsn := flag.String("d", "", "реквизиты базы данных")
	file := flag.String("f", "", "путь к файлу для хранения данных")
	trustedSubnet := flag.String("t", "", "доверенная подсеть")
	conf := flag.String("c", "", "Файл конфигурации")
	flag.StringVar(conf, "config", "", "Файл конфигурации")
	enableHTTPS := flag.Bool("s", false, "Enable HTTPS")

	flag.Parse()

	c.ServerAddress = *serverAddress
	c.DatabaseDSN = *dsn
	c.ConfigPath = *conf

	// С bool сложнее: флаг всегда false по умолчанию.
	// Проверяем, был ли он явно передан в командной строке.
	if isFlagPassed("s") {
		c.EnableHTTPS = *enableHTTPS
	}

	return c
}

func parseEnv() Config {
	return Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		DatabaseDSN:   os.Getenv("DATABASE_DSN"),
		ConfigPath:    os.Getenv("CONFIG"),
		EnableHTTPS:   os.Getenv("ENABLE_HTTPS") == "true",
	}
}

func parseJSON(path string) Config {
	var c Config

	if path == "" {
		return c
	}

	file, err := os.Open(path)

	if err != nil {
		return c
	}

	defer file.Close()
	json.NewDecoder(file).Decode(&c)

	return c
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})

	return found
}
