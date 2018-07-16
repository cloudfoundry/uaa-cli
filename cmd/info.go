package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func InfoCmd(api *uaa.API) error {
	i, err := api.GetInfo()
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
		NotifyErrorsWithRetry(InfoCmd(GetUnauthenticatedAPI()), log)
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
	infoCmd.Annotations = make(map[string]string)
	infoCmd.Annotations[INTRO_CATEGORY] = "true"
}
