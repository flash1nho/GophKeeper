package cmd

import (
	"github.com/spf13/cobra"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"

	"github.com/flash1nho/GophKeeper/app/client/cmd/secrets"
	"github.com/flash1nho/GophKeeper/config"
	"google.golang.org/grpc"
)

func SecretsCommand(settings config.SettingsObject) *cobra.Command {
	var client pb.GophKeeperPrivateServiceClient
	var conn *grpc.ClientConn

	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Менеджер хранения данных",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error

			token, err := cmd.Flags().GetString("token")

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			conn, err = BaseConnection(settings, token)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			client = pb.NewGophKeeperPrivateServiceClient(conn)

			if client == nil {
				settings.Log.Fatal("Ошибка при инициализации gRPC клиента")
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if err := conn.Close(); err != nil {
				settings.Log.Fatal(err.Error())
			}
		},
	}

	cmd.AddCommand(secrets.SecretsTextCommand(&client, settings))

	return cmd
}
