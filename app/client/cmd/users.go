package cmd

import (
	"github.com/spf13/cobra"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"

	"github.com/flash1nho/GophKeeper/app/client/cmd/users"
	"github.com/flash1nho/GophKeeper/config"
	"google.golang.org/grpc"
)

func UsersCommand(settings config.SettingsObject) *cobra.Command {
	var client pb.GophKeeperPublicServiceClient
	var conn *grpc.ClientConn

	cmd := &cobra.Command{
		Use:   "users",
		Short: "Менеджер регистрации и авторизации",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error

			conn, err = BaseConnection(settings, "")

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			client = pb.NewGophKeeperPublicServiceClient(conn)

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

	cmd.AddCommand(users.UsersRegisterCommand(&client, settings))
	cmd.AddCommand(users.UsersLoginCommand(&client, settings))

	return cmd
}
