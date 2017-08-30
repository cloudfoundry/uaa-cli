package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jhamon/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
)

var getClientCmd = &cobra.Command{
	Use:   "get-client CLIENT_ID",
	Short: "View client registration",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureTarget()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}
		client, err := cm.Get(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		j, err := json.MarshalIndent(&client, "", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(j))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return MissingArgument("client_id")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(getClientCmd)
}
