package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Регистрация нового пользователя в GophKeeper",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := &pb.UserRegisterRequest{
			Login:    login,
			Password: password,
		}

		client, conn, err := PublicClient()

		if err != nil {
			return err
		}

		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.Register(ctx, req)

		if err != nil {
			return fmt.Errorf("ошибка регистрации: %v", err)
		}

		fmt.Printf("Пользователь успешно зарегистрирован! Token: %s\n", resp.Token)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(&login, "login", "l", "", "login пользователя (обязательно)")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "login (обязательно)")

	registerCmd.MarkFlagRequired("login")
	registerCmd.MarkFlagRequired("password")
}
