package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func GetTokenKeysCmd(api *uaa.API) error {
	key, err := api.TokenKeys()

	if err != nil {
		return err
	}

	return cli.NewJsonPrinter(log).Print(key)
}

var getTokenKeysCmd = &cobra.Command{
	Use:     "get-token-keys",
	Short:   "View all keys the UAA has used to sign JWT tokens",
	Aliases: []string{"token-keys"},
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(cli.EnsureTargetInConfig(cfg), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli.NotifyErrorsWithRetry(GetTokenKeysCmd(GetUnauthenticatedAPI()), log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(getTokenKeysCmd)
	getTokenKeysCmd.Annotations = make(map[string]string)
	getTokenKeysCmd.Annotations[TOKEN_CATEGORY] = "true"
}
