package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
	"code.cloudfoundry.org/uaa-cli/cli"
)

var getClientCmd = &cobra.Command{
	Use:   "get-client CLIENT_ID",
	Short: "View client registration",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureContext()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}
		client, err := cm.Get(args[0])
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		err = cli.NewJsonPrinter(log).Print(client)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			MissingArgument("client_id", cmd)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(getClientCmd)
	getClientCmd.Annotations = make(map[string]string)
	getClientCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	getClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to get the client")
}
