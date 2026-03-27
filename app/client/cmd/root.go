package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/pkg/version"
)

var (
	buildVersion string = "N/A"
	buildCommit  string = "N/A"
	buildDate    string = "N/A"
)

var (
	login    string
	password string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gophkeeper",
	Short:   "GophKeeper client",
	Version: version.Info(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}

func init() {
	// TODO
}
