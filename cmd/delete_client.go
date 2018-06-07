package cmd

import (
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func DeleteClientValidations(cfg uaa.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) == 0 {
		return MissingArgumentError("client_id")
	}
	return nil
}

func DeleteClientCmd(cm *uaa.ClientManager, clientId string) error {
	_, err := cm.Delete(clientId)
	if err != nil {
		return err
	}

	log.Infof("Successfully deleted client %v.", utils.Emphasize(clientId))
	return nil
}

var deleteClientCmd = &cobra.Command{
	Use:   "delete-client CLIENT_ID",
	Short: "Delete a client registration",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(DeleteClientValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}
		NotifyErrorsWithRetry(DeleteClientCmd(cm, args[0]), GetSavedConfig(), log)
	},
}

func init() {
	RootCmd.AddCommand(deleteClientCmd)
	deleteClientCmd.Annotations = make(map[string]string)
	deleteClientCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	deleteClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to delete the client")
}
