package cmd

import (
	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var contextsCmd = cobra.Command{
	Use:   "contexts",
	Short: "List available contexts for the currently targeted UAA",
	Long:  help.Context(),
	Run: func(cmd *cobra.Command, args []string) {
		c := GetSavedConfig()

		if c.ActiveTargetName == "" {
			log.Error("No contexts are currently available.")
			log.Error(`To get started, target a UAA and fetch a token. See "uaa target -h" for details.`)
			os.Exit(1)
		}

		if len(c.GetActiveTarget().Contexts) == 0 {
			log.Error("No contexts are currently available.")
			log.Error(`Use a token command such as "uaa get-password-token" or "uaa get-client-credentials-token" to fetch a token and create a context.`)
			os.Exit(1)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ClientId", "Username", "Grant Type"})
		for _, context := range c.GetActiveTarget().Contexts {
			table.Append([]string{context.ClientId, context.Username, string(context.GrantType)})
		}
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(&contextsCmd)
	contextsCmd.Annotations = make(map[string]string)
	contextsCmd.Annotations[INTRO_CATEGORY] = "true"
	contextsCmd.Hidden = true
}
