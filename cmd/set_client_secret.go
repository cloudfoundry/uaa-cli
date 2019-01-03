package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func SetClientSecretValidation(cfg config.Config, args []string, clientSecret string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) == 0 {
		return MissingArgumentError("client_id")
	}
	if clientSecret == "" {
		return MissingArgumentError("client_secret")
	}
	return nil
}

func SetClientSecretCmd(api *uaa.API, log cli.Logger, clientId, clientSecret string) error {
	err := api.ChangeClientSecret(clientId, clientSecret)
	if err != nil {
		return errors.New("The secret for client " + clientId + " was not updated.")
	}
	log.Infof("The secret for client %v has been successfully updated.", utils.Emphasize(clientId))
	return nil
}

var setClientSecretCmd = &cobra.Command{
	Use:   "set-client-secret CLIENT_ID -s CLIENT_SECRET",
	Short: "Update secret for a client",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(SetClientSecretValidation(GetSavedConfig(), args, clientSecret), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		NotifyErrorsWithRetry(SetClientSecretCmd(GetAPIFromSavedTokenInContext(), log, args[0], clientSecret), log)
	},
}

func init() {
	RootCmd.AddCommand(setClientSecretCmd)
	setClientSecretCmd.Annotations = make(map[string]string)
	setClientSecretCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	setClientSecretCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "new client secret")
	setClientSecretCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain where the client resides")
}
