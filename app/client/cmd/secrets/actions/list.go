package actions

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

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

			for _, secretVal := range response.Secrets.Values {
				fields := secretVal.GetStructValue().GetFields()

				if id, ok := fields["id"]; ok {
					fmt.Printf("id: %v\n", id.GetNumberValue())
				}

				if data, ok := fields["data"]; ok {
					dataFields := data.GetStructValue().GetFields()
					keys := make([]string, 0, len(dataFields))

					for k := range dataFields {
						keys = append(keys, k)
					}

					sort.Strings(keys)

					for _, k := range keys {
						v := dataFields[k]

						fmt.Printf("%s: %s\n", k, v.GetStringValue())
					}
				}

				if createdAt, ok := fields["created_at"]; ok {
					fmt.Printf("created_at: %v\n", createdAt.GetStringValue())
				}

				if updatedAt, ok := fields["updated_at"]; ok {
					fmt.Printf("updated_at: %v\n", updatedAt.GetStringValue())
				}

				fmt.Println("---")
			}
		},
	}

	return cmd
}
