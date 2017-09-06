package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
)

var setClientSecretCmd = &cobra.Command{
	Use:   "set-client-secret CLIENT_ID -s CLIENT_SECRET",
	Short: "Update secret for a client",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureContext()
	},
	Run: func(cmd *cobra.Command, args []string) {
		clientId := args[0]
		c := GetSavedConfig()
		cm := &uaa.ClientManager{GetHttpClient(), c}
		err := cm.ChangeSecret(clientId, clientSecret)
		if err != nil {
			log.Errorf("The secret for client %v was not updated.", clientId)
			TraceRetryMsg(c)
			os.Exit(1)
		}
		log.Infof("The secret for client %v has been successfully updated.", clientId)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return MissingArgument("client_id")
		}
		if clientSecret == "" {
			return MissingArgument("client_secret")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(setClientSecretCmd)
	setClientSecretCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "new client secret")
	setClientSecretCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain where the client resides")
}
