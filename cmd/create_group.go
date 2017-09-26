package cmd

import (
	"errors"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
)

func CreateGroupCmd(gm uaa.GroupManager, printer cli.Printer, name string) error {
	toCreate := uaa.ScimGroup{
		DisplayName: name,
	}

	group, err := gm.Create(toCreate)
	if err != nil {
		return err
	}

	return printer.Print(group)
}

func CreateGroupValidation(cfg uaa.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("The positional argument GROUPNAME must be specified.")
	}
	return nil
}

var createGroupCmd = &cobra.Command{
	Use:     "create-group GROUPNAME",
	Short:   "Create a group",
	Aliases: []string{"add-group"},
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(CreateGroupValidation(GetSavedConfig(), args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		gm := uaa.GroupManager{GetHttpClient(), cfg}
		err := CreateGroupCmd(gm, cli.NewJsonPrinter(log), args[0])
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(createGroupCmd)
	createGroupCmd.Annotations = make(map[string]string)
	createGroupCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"
}
