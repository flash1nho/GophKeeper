package actions

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/iancoleman/strcase"
)

func SecretsCreateCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "create",
		Short:              "Создание секрета",
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

			Data, err := structpb.NewStruct(dataMap)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			Type := strcase.ToCamel(cmd.Parent().Name())

			request := &pb.CreateRequest{
				Data: Data,
				Type: Type,
			}

			response, err := (*client).Create(cmd.Context(), request)

			if err != nil {
				if statusErr, ok := status.FromError(err); ok {
					fmt.Printf("Ошибка создания секрета: %s\n", statusErr.Message())
				} else {
					settings.Log.Error(err.Error())
				}

				return
			}

			fmt.Printf("Создано! id: %d\n", int(response.ID))
		},
	}

	return cmd
}
