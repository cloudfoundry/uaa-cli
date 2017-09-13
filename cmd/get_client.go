package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"code.cloudfoundry.org/uaa-cli/cli"
)

func GetClientCmd(cm *uaa.ClientManager, clientId string) error {
	client, err := cm.Get(clientId)
	if err != nil {
		return err
	}

	return cli.NewJsonPrinter(log).Print(client)
}

func GetClientValidations(cfg uaa.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) == 0 {
		return MissingArgumentError("client_id")
	}
	return nil
}

var getClientCmd = &cobra.Command{
	Use:   "get-client CLIENT_ID",
	Short: "View client registration",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(GetClientValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}
		NotifyErrorsWithRetry(GetClientCmd(cm, args[0]), GetSavedConfig(), log)
	},
}

func init() {
	RootCmd.AddCommand(getClientCmd)
	getClientCmd.Annotations = make(map[string]string)
	getClientCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	getClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to get the client")
}
