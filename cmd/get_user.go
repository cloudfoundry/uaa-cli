package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
)

func GetUserCmd(userId string, um uaa.Crud, printer cli.Printer) error {
	user, err := um.Get(userId)
	if err != nil {
		return err
	}

	return printer.Print(user)
}

var getUserCmd = &cobra.Command{
	Use:   "get-user USER-ID",
	Short: "Look up a user by userId",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureContext()
	},
	Run: func(cmd *cobra.Command, args []string) {
		um := uaa.UserManager{GetHttpClient(), GetSavedConfig()}
		err := GetUserCmd(args[0], um, cli.NewJsonPrinter(log))
		if err != nil {
			log.Error(err.Error())
			TraceRetryMsg(GetSavedConfig())
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(getUserCmd)
	getUserCmd.Annotations = make(map[string]string)
	getUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"
}
