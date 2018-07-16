package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func GetTokenKeyCmd(api *uaa.API) error {
	key, err := api.TokenKey()
	if err != nil {
		return err
	}

	return cli.NewJsonPrinter(log).Print(key)
}

var getTokenKeyCmd = &cobra.Command{
	Use:     "get-token-key",
	Short:   "View the key for validating UAA's JWT token signatures",
	Aliases: []string{"token-key"},
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(EnsureTargetInConfig(cfg), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		NotifyErrorsWithRetry(GetTokenKeyCmd(GetUnauthenticatedAPI()), log)
	},
}

func init() {
	RootCmd.AddCommand(getTokenKeyCmd)
	getTokenKeyCmd.Annotations = make(map[string]string)
	getTokenKeyCmd.Annotations[TOKEN_CATEGORY] = "true"
}
