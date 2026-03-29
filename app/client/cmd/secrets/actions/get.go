package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"

	"github.com/iancoleman/strcase"
)

func SecretsGetCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	var ID int32

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Просмотр секрета",
		Run: func(cmd *cobra.Command, args []string) {
			Type := strcase.ToCamel(cmd.Parent().Name())

			request := &pb.GetRequest{
				ID:   ID,
				Type: Type,
			}

			response, err := (*client).Get(cmd.Context(), request)

			if err != nil {
				if statusErr, ok := status.FromError(err); ok {
					fmt.Printf("Ошибка просмотра секрета: %s\n", statusErr.Message())
				} else {
					settings.Log.Error(err.Error())
				}

				return
			}

			for key, value := range response.Secret.Fields {
				fmt.Printf("%s: %v\n", key, value.AsInterface())
			}
		},
	}

	cmd.Flags().Int32VarP(&ID, "id", "", 0, "ID (обязательно)")
	cmd.MarkFlagRequired("id")

	return cmd
}
