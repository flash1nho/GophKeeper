package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

func SecretsCreateCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject, fields []helpers.FieldInfo) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Создание секрета",
		Run: func(cmd *cobra.Command, args []string) {
			_, data, secretType, err := helpers.ArgsParse(cmd)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			request := &pb.CreateRequest{
				Data: data,
				Type: secretType,
			}

			if client == nil || *client == nil {
				settings.Log.Fatal("gRPC клиент не инициализирован")
			}

			response, err := (*client).Create(cmd.Context(), request)

			if err != nil {
				helpers.ErrorHandler(settings.Log, err)

				return
			}

			fmt.Println("✅ Создано")
			fmt.Println("---")

			if response != nil && response.Secrets != nil && response.Secrets.Values != nil {
				helpers.PrintResult(response.Secrets.Values)
			}
		},
	}

	for _, field := range fields {
		switch field.Type {
		case "string":
			var str string

			cmd.Flags().StringVarP(&str, field.Key, "", str, fmt.Sprintf("%s (обязательно)", field.Key))
			_ = cmd.MarkFlagRequired(field.Key)
		default:
			settings.Log.Fatal("недопустимый тип для создания флага")
		}
	}

	return cmd
}
