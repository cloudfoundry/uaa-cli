package cmd

import (
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
	"code.cloudfoundry.org/uaa-cli/cli"
)

var userinfoCmd = cobra.Command{
	Use:     "userinfo",
	Short:   "See claims about the authenticated user",
	Aliases: []string{"me"},
	Long:    help.Userinfo(),
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureContext()
	},
	Run: func(cmd *cobra.Command, args []string) {
		i, err := uaa.Me(GetHttpClient(), GetSavedConfig())
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
	RootCmd.AddCommand(&userinfoCmd)
	userinfoCmd.Annotations = make(map[string]string)
	userinfoCmd.Annotations[MISC_CATEGORY] = "true"
}
