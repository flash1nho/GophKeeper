package cmd

import (
	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:     "gophkeeper",
	Short:   "GophKeeper client",
	Version: version.Info(),
}

func Execute() {
	settings := config.Settings()

	rootCmd.AddCommand(UsersCommand(settings))
	rootCmd.AddCommand(SecretsCommand(settings))

	rootCmd.Execute()
}
