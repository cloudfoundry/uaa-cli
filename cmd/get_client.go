package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func GetClientCmd(api *uaa.API, clientId string) error {
	client, err := api.GetClient(clientId)
	if err != nil {
		return err
	}

	return cli.NewJsonPrinter(log).Print(client)
}

func GetClientValidations(cfg config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) == 0 {
		return cli.MissingArgumentError("client_id")
	}
	return nil
}

var getClientCmd = &cobra.Command{
	Use:   "get-client CLIENT_ID",
	Short: "View client registration",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(GetClientValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		api := NewApiFromSavedConfig()
		cli.NotifyErrorsWithRetry(GetClientCmd(api, args[0]), log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(getClientCmd)
	getClientCmd.Annotations = make(map[string]string)
	getClientCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	getClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to get the client")
}
