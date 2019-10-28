package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func MapGroupCmd(api *uaa.API, printer cli.Printer, externalGroupName string, groupName string) error {
	group, err := api.GetGroupByName(groupName, "")
	if err != nil {
		return err
	}
	err = api.MapGroup(group.ID, externalGroupName, "")
	if err != nil {
		return err
	}

	log.Infof("Successfully mapped %v to %v for origin %v", utils.Emphasize(groupName), utils.Emphasize(externalGroupName), utils.Emphasize("ldap"))
	return nil
}

func MapGroupValidations(cfg config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) != 2 {
		return errors.New("The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified.")
	}

	return nil
}

var mapGroupCmd = &cobra.Command{
	Use:   "map-group EXTERNAL_GROUPNAME GROUPNAME",
	Short: "Map uaa groups (scopes) to external groups defined within an external identity provider",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(MapGroupValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := MapGroupCmd(GetAPIFromSavedTokenInContext(), cli.NewJsonPrinter(log), args[0], args[1])
		cli.NotifyErrorsWithRetry(err, log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(mapGroupCmd)
	mapGroupCmd.Annotations = make(map[string]string)
	mapGroupCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"
}
