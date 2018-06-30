package cmd

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/spf13/cobra"
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
		if showAccessToken {
			fmt.Print(activeContext.Token.AccessToken)
			os.Exit(0)
		}
		if showAuthHeader {
			fmt.Printf(`%s %s`, activeContext.Token.TokenType, activeContext.Token.AccessToken)
			os.Exit(0)
		}
		err := cli.NewJsonPrinter(log).Print(activeContext)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(contextCmd)
	contextCmd.Annotations = make(map[string]string)
	contextCmd.Annotations[INTRO_CATEGORY] = "true"
	contextCmd.Flags().BoolVarP(&showAccessToken, "access_token", "", false, "Display context's access token")
	contextCmd.Flags().BoolVarP(&showAuthHeader, "auth_header", "a", false, "Display context's token type and access token")
}
