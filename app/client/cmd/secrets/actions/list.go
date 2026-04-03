package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

func SecretsListCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Список секретов",
		Run: func(cmd *cobra.Command, args []string) {
			_, _, secretType, err := helpers.ArgsParse(cmd)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			request := &pb.ListRequest{
				Type: secretType,
			}

			response, err := (*client).List(cmd.Context(), request)

			if err != nil {
				helpers.ErrorHandler(settings.Log, err)

				return
			}

			fmt.Println("Список секретов:")
			fmt.Println("---")

			if response != nil && response.Secrets != nil && response.Secrets.Values != nil {
				helpers.PrintResult(response.Secrets.Values)
			}
		},
	}

	return cmd
}
