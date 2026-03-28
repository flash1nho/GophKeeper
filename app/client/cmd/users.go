package cmd

import (
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"

	pb "github.com/flash1nho/GophKeeper/internal/grpc"

	"github.com/flash1nho/GophKeeper/app/client/cmd/users"
)

var (
	client   pb.GophKeeperPublicServiceClient
	login    string
	password string
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Менеджер регистрации и авторизации",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		conn, err := BaseConnection()

		if err != nil {
			log.Fatal().Err(err)
		}

		client = pb.NewGophKeeperPublicServiceClient(conn)
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	registerCmd := users.NewRegisterCmd(&client)
	registerCmd.Flags().StringVarP(&login, "login", "l", "", "Логин пользователя (обязательно)")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль пользователя (обязательно)")
	registerCmd.MarkFlagRequired("login")
	registerCmd.MarkFlagRequired("password")

	loginCmd := users.NewLoginCmd(&client)
	loginCmd.Flags().StringVarP(&login, "login", "l", "", "Логин пользователя (обязательно)")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль пользователя (обязательно)")
	loginCmd.MarkFlagRequired("login")
	loginCmd.MarkFlagRequired("password")

	usersCmd.AddCommand(registerCmd)
	usersCmd.AddCommand(loginCmd)
}
