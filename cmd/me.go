package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/jhamon/uaa-cli/uaa"
	"encoding/json"
	"github.com/jhamon/uaa-cli/help"
)

var meCmd = cobra.Command{
	Use:   "me",
	Short: "See claims about the authenticated user",
	Aliases: []string{"userinfo"},
	Long: help.Userinfo(),
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureTarget()
	},
	Run: func(cmd *cobra.Command, args []string) {
		i, err := uaa.Me(GetHttpClient(), GetSavedConfig())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		j, err := json.MarshalIndent(&i, "", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(j))
	},
}

func init() {
	RootCmd.AddCommand(&meCmd)
}
