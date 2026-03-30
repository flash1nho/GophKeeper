package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

func SecretsGetCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	var id int

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Просмотр секрета",
		Run: func(cmd *cobra.Command, args []string) {
			id, _, secretType, err := helpers.ArgsParse(cmd)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			request := &pb.GetRequest{
				ID:   int32(id),
				Type: secretType,
			}

			response, err := (*client).Get(cmd.Context(), request)

			if err != nil {
				helpers.ErrorHandler(settings.Log, err)

				return
			}

			fmt.Println("Карточка секрета:")
			fmt.Println("---")

			if response != nil && response.Secrets != nil && response.Secrets.Values != nil {
				helpers.PrintResult(response.Secrets.Values)
			}
		},
	}

	cmd.Flags().IntVarP(&id, "id", "", id, "id (обязательно)")
	cmd.MarkFlagRequired("id")

	return cmd
}
