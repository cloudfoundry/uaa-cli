package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
	"code.cloudfoundry.org/uaa-cli/cli"
)

var listClientsCmd = &cobra.Command{
	Use:     "list-clients",
	Short:   "See all clients in the targeted UAA",
	Aliases: []string{"clients"},
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureContext()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}
		clients, err := cm.List()
		if err != nil {
			log.Error(err.Error())
			TraceRetryMsg(GetSavedConfig())
			os.Exit(1)
		}

		err = cli.NewJsonPrinter(log).Print(clients)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(listClientsCmd)
	listClientsCmd.Annotations = make(map[string]string)
	listClientsCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	listClientsCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to get the client")
}
