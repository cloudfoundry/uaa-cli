package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

var oldSecret string

func ChangeClientSecretValidation(cfg config.Config, oldSecret, newSecret string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}

	context := cfg.GetActiveContext()
	if context.GrantType != config.CLIENT_CREDENTIALS {
		return errors.New("You must have a client_credentials token in your context to perform this command.")
	}

	if oldSecret == "" {
		return cli.MissingArgumentError("old_secret")
	}
	if newSecret == "" {
		return cli.MissingArgumentError("secret")
	}
	return nil
}

func ChangeClientSecretCmd(api *uaa.API, log cli.Logger, cfg config.Config, oldSecret, newSecret string) error {
	context := cfg.GetActiveContext()
	clientId := context.ClientId

	// Prepare the request body for the secret change
	requestBody := map[string]interface{}{
		"oldSecret": oldSecret,
		"secret":    newSecret,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// Make the API call to change the client secret
	path := fmt.Sprintf("/oauth/clients/%s/secret", clientId)
	headers := []string{"Content-Type: application/json"}

	// Add zone header if specified
	if cfg.ZoneSubdomain != "" {
		headers = append(headers, fmt.Sprintf("X-Identity-Zone-Id: %s", cfg.ZoneSubdomain))
	}

	_, _, status, err := api.Curl(path, "PUT", string(requestBodyJSON), headers)
	if err != nil {
		return err
	}

	if status >= 400 {
		return errors.New("The secret for client " + clientId + " was not updated.")
	}

	log.Infof("The secret for client %v has been successfully updated.", utils.Emphasize(clientId))
	return nil
}

var changeClientSecretCmd = &cobra.Command{
	Use:   "change-client-secret --old_secret OLD_SECRET --secret NEW_SECRET",
	Short: "Change secret for authenticated client",
	PreRun: func(cmd *cobra.Command, args []string) {
		cli.NotifyValidationErrors(ChangeClientSecretValidation(GetSavedConfig(), oldSecret, clientSecret), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		api := GetAPIFromSavedTokenInContext()
		cli.NotifyErrorsWithRetry(ChangeClientSecretCmd(api, log, cfg, oldSecret, clientSecret), log, cfg)
	},
}

func init() {
	RootCmd.AddCommand(changeClientSecretCmd)
	changeClientSecretCmd.Annotations = make(map[string]string)
	changeClientSecretCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	changeClientSecretCmd.Flags().StringVar(&oldSecret, "old_secret", "", "current client secret")
	changeClientSecretCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "new client secret")
	changeClientSecretCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain where the client resides")
}
