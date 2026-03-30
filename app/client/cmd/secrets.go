package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/flash1nho/GophKeeper/app/client/cmd/secrets"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

func SecretsCommand(settings config.SettingsObject) *cobra.Command {
	var client pb.GophKeeperPrivateServiceClient
	var conn *grpc.ClientConn

	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Менеджер хранения данных",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			token := viper.GetString("token")

			if token == "" {
				settings.Log.Fatal("Токен не найден. Сначала выполните вход (login)!")
			}

			conn, err := BaseConnection(settings, token)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			client = pb.NewGophKeeperPrivateServiceClient(conn)

			if client == nil {
				settings.Log.Fatal("Не удалось инициализировать gRPC клиента")
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if conn != nil {
				if err := conn.Close(); err != nil {
					settings.Log.Fatal(err.Error())
				}
			}
		},
	}

	cmd.AddCommand(secrets.TextCommand(&client, settings))
	cmd.AddCommand(secrets.CredCommand(&client, settings))
	cmd.AddCommand(secrets.CardCommand(&client, settings))

	return cmd
}
