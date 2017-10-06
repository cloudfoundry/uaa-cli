package cmd

import (
	"github.com/spf13/cobra"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"net/http"
	"code.cloudfoundry.org/uaa-cli/cli"
)

func GetTokenKeyValidations(cfg uaa.Config) error {
	if err := EnsureTargetInConfig(cfg); err != nil {
		return err
	}
	return nil
}

func GetTokenKeyCmd(client *http.Client, config uaa.Config) error {
	key, err := uaa.TokenKey(client, config)

	if err != nil {
		return err
	}

	return cli.NewJsonPrinter(log).Print(key)
}

var getTokenKeyCmd = &cobra.Command{
	Use:   "get-token-key",
	Short: "View the key for validating UAA's JWT token signatures",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(GetTokenKeyValidations(cfg), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		NotifyErrorsWithRetry(GetTokenKeyCmd(GetHttpClient(), GetSavedConfig()), GetSavedConfig(), log)
	},
}

func init() {
	RootCmd.AddCommand(getTokenKeyCmd)
	getTokenKeyCmd.Annotations = make(map[string]string)
	getTokenKeyCmd.Annotations[TOKEN_CATEGORY] = "true"
}
