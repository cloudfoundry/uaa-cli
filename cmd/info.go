package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"encoding/json"
	"github.com/spf13/cobra"
	"os"
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

		j, err := json.MarshalIndent(&i, "", "  ")
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		log.Robots(string(j))
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
	infoCmd.Annotations = make(map[string]string)
	infoCmd.Annotations[INTRO_CATEGORY] = "true"
}
