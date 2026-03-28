package users

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"
)

func NewLoginCmd(client *pb.GophKeeperPublicServiceClient) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Авторизация пользователя",
		Run: func(cmd *cobra.Command, args []string) {
			if client == nil || *client == nil {
				log.Fatal().Err(fmt.Errorf("gRPC клиент не инициализирован"))

				return
			}

			login, err := cmd.Flags().GetString("login")

			if err != nil {
				log.Fatal().Err(err)

				return
			}

			password, err := cmd.Flags().GetString("password")

			if err != nil {
				log.Fatal().Err(err)

				return
			}

			request := &pb.UserLoginRequest{
				Login:    login,
				Password: password,
			}

			response, err := (*client).Login(cmd.Context(), request)

			if err != nil {
				if s, ok := status.FromError(err); ok {
					fmt.Printf("Ошибка авторизации: %s\n", s.Message())
				} else {
					fmt.Printf("Неизвестная ошибка: %v", err)
				}

				return
			}

			fmt.Printf("Пользователь успешно авторизован!\n\nToken: %s\n", response.Token)
		},
	}
}
