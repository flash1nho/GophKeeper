package actions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"

	"github.com/iancoleman/strcase"
)

func SecretsGetCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	id := 0

	cmd := &cobra.Command{
		Use:                "get",
		Short:              "Просмотр секрета",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			dataMap := make(map[string]interface{})

			for _, arg := range args {
				kv := strings.SplitN(arg, "=", 2)

				if len(kv) == 2 {
					field := strings.TrimPrefix(kv[0], "--")
					dataMap[field] = kv[1]
				}
			}

			idStr, ok := dataMap["id"].(string)

			if ok {
				fmt.Sscanf(idStr, "%d", &id)
			}

			Type := strcase.ToCamel(cmd.Parent().Name())

			request := &pb.GetRequest{
				ID:   int32(id),
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

			fields := response.Secret.GetFields()

			if id, ok := fields["id"]; ok {
				fmt.Printf("id: %g\n", id.GetNumberValue())
			}

			if data, ok := fields["data"]; ok {
				if structData := data.GetStructValue(); structData != nil {
					dataFields := structData.GetFields()
					keys := make([]string, 0, len(dataFields))

					for k := range dataFields {
						keys = append(keys, k)
					}

					sort.Strings(keys)

					for _, k := range keys {
						v := dataFields[k]

						fmt.Printf("%s: %v\n", k, v.AsInterface())
					}
				} else {
					fmt.Printf("data: %v\n", data.AsInterface())
				}
			}

			if createdAt, ok := fields["created_at"]; ok {
				fmt.Printf("created_at: %v\n", createdAt.GetStringValue())
			}

			if updatedAt, ok := fields["updated_at"]; ok {
				fmt.Printf("updated_at: %v\n", updatedAt.GetStringValue())
			}
		},
	}

	return cmd
}
