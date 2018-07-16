package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func ListClientsValidations(cfg config.Config) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	return nil
}

func ListClientsCmd(api *uaa.API) error {
	clients, err := api.ListAllClients("", "", uaa.SortAscending)
	if err != nil {
		return err
	}

	return cli.NewJsonPrinter(log).Print(clients)
}

var listClientsCmd = &cobra.Command{
	Use:     "list-clients",
	Short:   "See all clients in the targeted UAA",
	Aliases: []string{"clients"},
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(ListClientsValidations(GetSavedConfig()), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		NotifyErrorsWithRetry(ListClientsCmd(GetAPIFromSavedTokenInContext()), log)
	},
}

func init() {
	RootCmd.AddCommand(listClientsCmd)
	listClientsCmd.Annotations = make(map[string]string)
	listClientsCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	listClientsCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to get the client")
}
