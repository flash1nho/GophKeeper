package actions

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"

	"github.com/iancoleman/strcase"
)

func SecretsDeleteCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	id := 0

	cmd := &cobra.Command{
		Use:                "delete",
		Short:              "Удаление секрета",
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

			request := &pb.DeleteRequest{
				ID:   int32(id),
				Type: Type,
			}

			_, err := (*client).Delete(cmd.Context(), request)

			if err != nil {
				if statusErr, ok := status.FromError(err); ok {
					fmt.Printf("Ошибка просмотра секрета: %s\n", statusErr.Message())
				} else {
					settings.Log.Error(err.Error())
				}

				return
			}

			fmt.Printf("Секрет с ID=%d удален!\n", id)
		},
	}

	return cmd
}
