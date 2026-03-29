package secrets

import (
	"github.com/spf13/cobra"

	"github.com/flash1nho/GophKeeper/app/client/cmd/secrets/actions"
	"github.com/flash1nho/GophKeeper/config"
	pb "github.com/flash1nho/GophKeeper/internal/grpc"
)

func CredCommand(client *pb.GophKeeperPrivateServiceClient, settings config.SettingsObject) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cred",
		Short: "Пара логин/пароль",
	}

	cmd.AddCommand(actions.SecretsCreateCommand(client, settings))
	cmd.AddCommand(actions.SecretsGetCommand(client, settings))
	cmd.AddCommand(actions.SecretsListCommand(client, settings))
	cmd.AddCommand(actions.SecretsUpdateCommand(client, settings))

	return cmd
}
