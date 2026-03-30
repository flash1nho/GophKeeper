package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

func SecretsUpdateCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject, fields []helpers.FieldInfo) *cobra.Command {
	var id int

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Обновление секрета",
		Run: func(cmd *cobra.Command, args []string) {
			id, data, secretType, err := helpers.ArgsParse(cmd)

			request := &pb.UpdateRequest{
				ID:   int32(id),
				Data: data,
				Type: secretType,
			}

			response, err := (*client).Update(cmd.Context(), request)

			if err != nil {
				helpers.ErrorHandler(settings.Log, err)

				return
			}

			fmt.Println("✅ Обновлено")
			fmt.Println("---")

			if response != nil && response.Secrets != nil && response.Secrets.Values != nil {
				helpers.PrintResult(response.Secrets.Values)
			}
		},
	}

	cmd.Flags().IntVarP(&id, "id", "", id, "id (обязательно)")
	cmd.MarkFlagRequired("id")

	for _, field := range fields {
		switch field.Type {
		case "string":
			var str string

			cmd.Flags().StringVarP(&str, field.Key, "", str, field.Key)
		default:
			settings.Log.Fatal("недопустимый тип для создания флага")
		}
	}

	return cmd
}
