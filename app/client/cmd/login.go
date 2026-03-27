package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Авторизация пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		req := &pb.UserLoginRequest{
			Login:    login,
			Password: password,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.Login(ctx, req)

		if err != nil {
			if s, ok := status.FromError(err); ok {
				fmt.Printf("Ошибка авторизации: %s\n", s.Message())
			} else {
				fmt.Printf("Неизвестная ошибка: %v", err)
			}

			return
		}

		fmt.Printf("Пользователь успешно авторизован!\n\nToken: %s\n", resp.Token)
	},
}

func init() {
	usersCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&login, "login", "l", "", "login пользователя (обязательно)")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "login (обязательно)")

	loginCmd.MarkFlagRequired("login")
	loginCmd.MarkFlagRequired("password")
}
