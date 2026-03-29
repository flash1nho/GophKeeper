package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:     "gophkeeper",
	Short:   "GophKeeper client",
	Version: version.Info(),
}

func Execute() {
	cobra.OnInitialize(initConfig)

	settings := config.Settings()

	rootCmd.AddCommand(UsersCommand(settings))
	rootCmd.AddCommand(SecretsCommand(settings))

	rootCmd.Execute()
}

func initConfig() {
	home, _ := os.UserHomeDir()
	cfgPath := filepath.Join(home, ".gophkeeper.yaml")

	viper.SetConfigFile(cfgPath)
	viper.SetConfigType("yaml")

	viper.ReadInConfig()
}
