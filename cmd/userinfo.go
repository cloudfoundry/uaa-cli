package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func UserinfoValidations(cfg config.Config) error {
	return EnsureContextInConfig(cfg)
}

func UserinfoCmd(api *uaa.API) error {
	i, err := api.GetMe()
	if err != nil {
		return err
	}

	return cli.NewJsonPrinter(log).Print(i)
}

var userinfoCmd = cobra.Command{
	Use:     "userinfo",
	Short:   "See claims about the authenticated user",
	Aliases: []string{"me"},
	Long:    help.Userinfo(),
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(UserinfoValidations(GetSavedConfig()), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := UserinfoCmd(GetAPIFromSavedTokenInContext())
		NotifyErrorsWithRetry(err, log)
	},
}

func init() {
	RootCmd.AddCommand(&userinfoCmd)
	userinfoCmd.Annotations = make(map[string]string)
	userinfoCmd.Annotations[MISC_CATEGORY] = "true"
}
