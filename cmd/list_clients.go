package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"encoding/json"
	"github.com/spf13/cobra"
	"os"
)

var listClientsCmd = &cobra.Command{
	Use:   "list-clients",
	Short: "See all clients in the targeted UAA",
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

		j, err := json.MarshalIndent(&clients, "", "  ")
		if err != nil {
			log.Error(err.Error())
			TraceRetryMsg(GetSavedConfig())
			os.Exit(1)
		}
		log.Robots(string(j))
	},
}

func init() {
	RootCmd.AddCommand(listClientsCmd)
	listClientsCmd.Annotations = make(map[string]string)
	listClientsCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	listClientsCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to get the client")
}
