package cmd

import (
	"errors"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func CreateGroupCmd(api *uaa.API, printer cli.Printer, name, description string) error {
	toCreate := uaa.Group{
		DisplayName: name,
		Description: description,
	}

	group, err := api.CreateGroup(toCreate)
	if err != nil {
		return err
	}

	return printer.Print(group)
}

func CreateGroupValidation(cfg config.Config, args []string) error {
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
		api, err := uaa.NewWithToken(cfg.GetActiveTarget().BaseUrl, cfg.ZoneSubdomain, cfg.GetActiveContext().Token, cfg.GetActiveTarget().SkipSSLValidation)
		NotifyErrorsWithRetry(err, log)
		err = CreateGroupCmd(api, cli.NewJsonPrinter(log), args[0], groupDescription)
		NotifyErrorsWithRetry(err, log)
	},
}

func init() {
	RootCmd.AddCommand(createGroupCmd)
	createGroupCmd.Annotations = make(map[string]string)
	createGroupCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"

	createGroupCmd.Flags().StringVarP(&groupDescription, "description", "d", "", `a human-readable description`)
	createGroupCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to create the group")
}
