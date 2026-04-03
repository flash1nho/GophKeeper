package secrets

import (
	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/cmd/secrets/actions"
	"github.com/flash1nho/GophKeeper/app/client/helpers"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
	"github.com/flash1nho/GophKeeper/internal/models/secrets"
)

func TextCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "text",
		Short: "Текстовые данные",
	}

	fields := helpers.GetStructKeys(secrets.Text{})

	cmd.AddCommand(actions.SecretsCreateCommand(client, settings, fields))
	cmd.AddCommand(actions.SecretsGetCommand(client, settings))
	cmd.AddCommand(actions.SecretsListCommand(client, settings))
	cmd.AddCommand(actions.SecretsUpdateCommand(client, settings, fields))
	cmd.AddCommand(actions.SecretsDeleteCommand(client, settings))

	return cmd
}
