package cmd

import (
	"context"
	"fmt"
	"time"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"

	"github.com/flash1nho/GophKeeper/internal/config"

	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

var (
	login    string
	password string
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Регистрация нового пользователя в GophKeeper",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(config.GrpcServerAddress, grpc.WithInsecure())

		if err != nil {
			return fmt.Errorf("не удалось подключиться: %v", err)
		}

		defer conn.Close()

		client := pb.NewGophKeeperServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		req := &pb.UserRegisterRequest{
			Login:    login,
			Password: password,
		}

		resp, err := client.Register(ctx, req)

		if err != nil {
			return fmt.Errorf("ошибка регистрации: %v", err)
		}

		fmt.Printf("Пользователь успешно зарегистрирован! ID: %s\n", resp.UserID)

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
