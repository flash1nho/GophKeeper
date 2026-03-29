package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

			token := viper.GetString("token")

			if token == "" {
				fmt.Println("Токен не найден. Выполните вход!")

				return
			}

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

	cmd.AddCommand(secrets.TextCommand(&client, settings))
	cmd.AddCommand(secrets.CredCommand(&client, settings))

	return cmd
}
