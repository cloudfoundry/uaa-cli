package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"code.cloudfoundry.org/uaa-cli/version"
)


var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.VersionString())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Annotations = make(map[string]string)
	versionCmd.Annotations[INTRO_CATEGORY] = "true"
}

