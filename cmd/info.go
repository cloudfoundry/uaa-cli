package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
	"code.cloudfoundry.org/uaa-cli/cli"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "See version and global configurations for the targeted UAA",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureTarget()
	},
	Run: func(cmd *cobra.Command, args []string) {
		i, err := uaa.Info(GetHttpClient(), GetSavedConfig())
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		err = cli.NewJsonPrinter(log).Print(i)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
	infoCmd.Annotations = make(map[string]string)
	infoCmd.Annotations[INTRO_CATEGORY] = "true"
}
