package cmd

import (
	"errors"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
)

func GetGroupCmd(gm uaa.GroupManager, printer cli.Printer, name, attributes string) error {
	group, err := gm.GetByName(name, attributes)
	if err != nil {
		return err
	}

	return printer.Print(group)
}

func GetGroupValidations(cfg uaa.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("The positional argument GROUPNAME must be specified.")
	}
	return nil
}

var getGroupCmd = &cobra.Command{
	Use:   "get-group GROUPNAME",
	Short: "Look up a group by groupname",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(GetGroupValidations(GetSavedConfig(), args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		gm := uaa.GroupManager{GetHttpClient(), cfg}
		err := GetGroupCmd(gm, cli.NewJsonPrinter(log), args[0], attributes)
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(getGroupCmd)
	getGroupCmd.Annotations = make(map[string]string)
	getGroupCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"

	getGroupCmd.Flags().StringVarP(&attributes, "attributes", "a", "", `include only these comma-separated user attributes to improve query performance`)
}
