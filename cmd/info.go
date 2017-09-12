package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"code.cloudfoundry.org/uaa-cli/cli"
	"net/http"
)

func InfoCmd(cfg uaa.Config, httpClient *http.Client) error {
	i, err := uaa.Info(httpClient, cfg)
	if err != nil {
		return err
	}

	return cli.NewJsonPrinter(log).Print(i)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "See version and global configurations for the targeted UAA",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(EnsureTargetInConfig(cfg), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyErrorsWithRetry(InfoCmd(cfg, GetHttpClient()), cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
	infoCmd.Annotations = make(map[string]string)
	infoCmd.Annotations[INTRO_CATEGORY] = "true"
}
