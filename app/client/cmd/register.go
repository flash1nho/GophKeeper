package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Регистрация пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		req := &pb.UserRegisterRequest{
			Login:    login,
			Password: password,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.Register(ctx, req)

		if err != nil {
			if s, ok := status.FromError(err); ok {
				fmt.Printf("Ошибка регистрации: %s\n", s.Message())
			} else {
				fmt.Printf("Неизвестная ошибка: %v", err)
			}

			return
		}

		fmt.Printf("Пользователь успешно зарегистрирован!\n\nToken: %s\n", resp.Token)
	},
}

func init() {
	usersCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(&login, "login", "l", "", "login пользователя (обязательно)")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "login (обязательно)")

	registerCmd.MarkFlagRequired("login")
	registerCmd.MarkFlagRequired("password")
}
