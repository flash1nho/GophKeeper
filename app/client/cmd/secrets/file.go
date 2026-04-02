package secrets

import (
	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/cmd/secrets/actions"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

var stateFile = ".upload_offsets"

func FileCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "Файлы",
	}

	cmd.AddCommand(actions.SecretsUploadCommand(client, settings))
	cmd.AddCommand(actions.SecretsDownloadCommand(client, settings))
	cmd.AddCommand(actions.SecretsGetCommand(client, settings))
	cmd.AddCommand(actions.SecretsListCommand(client, settings))
	cmd.AddCommand(actions.SecretsDeleteCommand(client, settings))

	return cmd
}
