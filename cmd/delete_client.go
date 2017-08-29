package cmd

import (
	"fmt"
	"os"
	"github.com/jhamon/uaa-cli/uaa"
	"github.com/spf13/cobra"
)

var deleteClientCmd = &cobra.Command{
	Use:   "delete-client CLIENT_ID",
	Short: "Delete a client registration",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureTarget()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}
		_, err := cm.Delete(args[0])
		if err != nil {
			fmt.Println(err)
			TraceRetryMsg(GetSavedConfig())
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted client %v.\n", args[0])
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return MissingArgument("client_id")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteClientCmd)
}
