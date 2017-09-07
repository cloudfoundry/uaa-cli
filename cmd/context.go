package cmd

import (
	"code.cloudfoundry.org/uaa-cli/help"
	"encoding/json"
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
			log.Error("No context is currently set.")
			log.Error(`To get started, target a UAA and fetch a token. See "uaa target -h" for details.`)
			os.Exit(1)
		}

		if len(c.GetActiveTarget().Contexts) == 0 {
			log.Error("No context is currently set.")
			log.Error(`Use a token command such as "uaa get-password-token" or "uaa get-client-credentials-token" to fetch a token.`)
			os.Exit(1)
		}

		activeContext := c.GetActiveContext()
		j, err := json.MarshalIndent(&activeContext, "", "  ")
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		log.Robots(string(j))
	},
}

func init() {
	RootCmd.AddCommand(contextCmd)
	contextCmd.Annotations = make(map[string]string)
	contextCmd.Annotations[INTRO_CATEGORY] = "true"
}
