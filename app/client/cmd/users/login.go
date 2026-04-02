package users

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
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
				helpers.ErrorHandler(settings.Log, err)

				return
			}

			viper.Set("token", response.Token)
			_ = viper.WriteConfig()
			_ = viper.SafeWriteConfig()

			fmt.Println("✅ Успешный вход!")
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Логин (обязательно)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Пароль (обязательно)")
	_ = cmd.MarkFlagRequired("login")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
