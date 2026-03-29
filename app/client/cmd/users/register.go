package users

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"
)

func UsersRegisterCommand(client *pb.GophKeeperPublicServiceClient, settings config.SettingsObject) *cobra.Command {
	var login string
	var password string
	var secret string

	cmd := &cobra.Command{
		Use:   "register",
		Short: "Регистрация пользователя",
		Run: func(cmd *cobra.Command, args []string) {
			request := &pb.UserRegisterRequest{
				Login:    login,
				Password: password,
				Secret:   secret,
			}

			response, err := (*client).Register(cmd.Context(), request)

			if err != nil {
				if statusErr, ok := status.FromError(err); ok {
					fmt.Printf("Ошибка регистрации: %s\n", statusErr.Message())
				} else {
					settings.Log.Error(err.Error())
				}

				return
			}

			viper.Set("token", response.Token)

			if err := viper.WriteConfig(); err != nil {
				viper.SafeWriteConfig()
			}

			fmt.Println("Успешная регистрация!")
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Логин (обязательно)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Пароль (обязательно)")
	cmd.Flags().StringVarP(&secret, "secret", "s", "", "Секретное слово (обязательно)")
	cmd.MarkFlagRequired("login")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("secret")

	return cmd
}
