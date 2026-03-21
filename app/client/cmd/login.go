package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Авторизация пользователя в GophKeeper",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := &pb.UserLoginRequest{
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

		resp, err := client.Login(ctx, req)

		if err != nil {
			return fmt.Errorf("ошибка авторизации: %v", err)
		}

		fmt.Printf("Пользователь успешно авторизован! Token: %s\n", resp.Token)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&login, "login", "l", "", "login пользователя (обязательно)")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "login (обязательно)")

	loginCmd.MarkFlagRequired("login")
	loginCmd.MarkFlagRequired("password")
}
