package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jhamon/uaa-cli/help"
	"github.com/spf13/cobra"
	"os"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "See information about the currently active CLI context",
	Long:  help.Context(),
	Run: func(cmd *cobra.Command, args []string) {
		c := GetSavedConfig()

		if c.ActiveTargetName == "" {
			fmt.Println("No context is currently set.")
			fmt.Println(`To get started, target a UAA and fetch a token. See "uaa target -h" for details.`)
			os.Exit(1)
		}

		if len(c.GetActiveTarget().Contexts) == 0 {
			fmt.Println("No context is currently set.")
			fmt.Println(`Use a token command such as "uaa get-password-token" or "uaa get-client-credentials-token" to fetch a token.`)
			os.Exit(1)
		}

		activeContext := c.GetActiveContext()
		j, err := json.MarshalIndent(&activeContext, "", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(j))

	},
}

func init() {
	RootCmd.AddCommand(contextCmd)
}
