package text

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"
)

func SecretsTextCreateCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	var token string
	var content string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Создание текстовых данных",
		Run: func(cmd *cobra.Command, args []string) {
			request := &pb.TextCreateRequest{
				Content: content,
			}

			response, err := (*client).TextCreate(cmd.Context(), request)

			if err != nil {
				if statusErr, ok := status.FromError(err); ok {
					fmt.Printf("Ошибка создания текстовых данных: %s\n", statusErr.Message())
				} else {
					settings.Log.Error(err.Error())
				}

				return
			}

			fmt.Printf("id: %d\n", int(response.ID))
		},
	}

	cmd.Flags().StringVarP(&content, "content", "p", "", "Содержимое (обязательно)")
	cmd.Flags().StringVarP(&token, "token", "t", "", "Token (обязательно)")
	cmd.MarkFlagRequired("content")
	cmd.MarkFlagRequired("token")

	return cmd
}
