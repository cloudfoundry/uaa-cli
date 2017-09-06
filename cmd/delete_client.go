package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
)

var deleteClientCmd = &cobra.Command{
	Use:   "delete-client CLIENT_ID",
	Short: "Delete a client registration",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureContext()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}
		_, err := cm.Delete(args[0])
		if err != nil {
			log.Error(err.Error())
			TraceRetryMsg(GetSavedConfig())
			os.Exit(1)
		}

		log.Infof("Successfully deleted client %v.", args[0])
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			MissingArgument("client_id", cmd)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteClientCmd)
	deleteClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to delete the client")
}
