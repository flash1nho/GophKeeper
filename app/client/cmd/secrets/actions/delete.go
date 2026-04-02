package actions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

func SecretsDeleteCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	var id int

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Удаление секрета",
		Run: func(cmd *cobra.Command, args []string) {
			id, _, secretType, err := helpers.ArgsParse(cmd)

			if err != nil {
				settings.Log.Fatal(err.Error())
			}

			request := &pb.DeleteRequest{
				ID:   int32(id),
				Type: secretType,
			}

			_, err = (*client).Delete(cmd.Context(), request)

			if err != nil {
				helpers.ErrorHandler(settings.Log, err)

				return
			}

			fmt.Printf("✅ Секрет с id=%d успешно удален!\n", id)
		},
	}

	cmd.Flags().IntVarP(&id, "id", "", id, "id (обязательно)")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
