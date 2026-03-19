package cmd

import (
	"github.com/flash1nho/GophKeeper/pkg/version"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gophkeeper-cli",
	Short:   "GophKeeper client",
	Version: version.Info(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("grpc-server-address", "g", "", "gRPC server address")
}
