package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"errors"
	"github.com/spf13/cobra"
)

func GetUserCmd(userId string, um uaa.Crud, printer cli.Printer) error {
	user, err := um.Get(userId)
	if err != nil {
		return err
	}

	return printer.Print(user)
}

func GetUserValidations(cfg uaa.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("The positional argument USER_ID must be specified.")
	}
	return nil
}

var getUserCmd = &cobra.Command{
	Use:   "get-user USER_ID",
	Short: "Look up a user by userId",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(GetUserValidations(GetSavedConfig(), args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		um := uaa.UserManager{GetHttpClient(), cfg}
		err := GetUserCmd(args[0], um, cli.NewJsonPrinter(log))
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(getUserCmd)
	getUserCmd.Annotations = make(map[string]string)
	getUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"
}
