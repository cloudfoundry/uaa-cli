package cmd

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func UserinfoValidations(cfg uaa.Config) error {
	return EnsureContextInConfig(cfg)
}

func UserinfoCmd(client *http.Client, cfg uaa.Config, printer cli.Printer) error {
	i, err := uaa.Me(GetHttpClient(), GetSavedConfig())
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
		cfg := GetSavedConfig()
		printer := cli.NewJsonPrinter(log)
		err := UserinfoCmd(GetHttpClient(), cfg, printer)
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(&userinfoCmd)
	userinfoCmd.Annotations = make(map[string]string)
	userinfoCmd.Annotations[MISC_CATEGORY] = "true"
}
