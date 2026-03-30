package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/cmd/secrets/print"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"

	"github.com/iancoleman/strcase"
)

func SecretsListCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "list",
		Short:              "Список секретов",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			Type := strcase.ToCamel(cmd.Parent().Name())

			request := &pb.ListRequest{
				Type: Type,
			}

			response, err := (*client).List(cmd.Context(), request)

			if err != nil {
				if statusErr, ok := status.FromError(err); ok {
					fmt.Printf("Ошибка списка секретов: %s\n", statusErr.Message())
				} else {
					settings.Log.Error(err.Error())
				}

				return
			}

			fmt.Println("Список секретов:")
			fmt.Println("---")
			print.Result(response.Secrets.Values)
		},
	}

	return cmd
}
