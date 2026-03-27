package cmd

import (
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc"
)

var (
	conn   grpc.ClientConn
	client pb.GophKeeperPublicServiceClient
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Менеджер регистрации и авторизации",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		conn, err := BaseConnection()

		if err != nil {
			log.Fatal().Err(err)
		}

		client = pb.NewGophKeeperPublicServiceClient(conn)

		if err != nil {
			log.Fatal().Err(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
