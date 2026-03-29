package users

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"
)

func UsersLoginCommand(client *pb.GophKeeperPublicServiceClient, settings config.SettingsObject) *cobra.Command {
	var login string
	var password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Авторизация пользователя",
		Run: func(cmd *cobra.Command, args []string) {
			request := &pb.UserLoginRequest{
				Login:    login,
				Password: password,
			}

			response, err := (*client).Login(cmd.Context(), request)

			if err != nil {
				if statusErr, ok := status.FromError(err); ok {
					fmt.Printf("Ошибка авторизации: %s\n", statusErr.Message())
				} else {
					settings.Log.Error(err.Error())
				}

				return
			}

			fmt.Printf("token: %s\n", response.Token)
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Логин (обязательно)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Пароль (обязательно)")
	cmd.MarkFlagRequired("login")
	cmd.MarkFlagRequired("password")

	return cmd
}
