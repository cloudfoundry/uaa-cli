package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func MapGroupCmd(api *uaa.API, printer cli.Printer, externalGroupName, groupName, origin string) error {
	if origin == "" {
		origin = "ldap"
	}
	group, err := api.GetGroupByName(groupName, "")
	if err != nil {
		return err
	}
	err = api.MapGroup(group.ID, externalGroupName, origin)
	if err != nil {
		return err
	}

	log.Infof("Successfully mapped %v to %v for origin %v", utils.Emphasize(groupName), utils.Emphasize(externalGroupName), utils.Emphasize(origin))
	return nil
}

var mapGroupCmd = &cobra.Command{
	Use:   "map-group EXTERNAL_GROUPNAME GROUPNAME",
	Short: "Map uaa groups (scopes) to external groups defined within an external identity provider",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(GroupMappingValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := MapGroupCmd(GetAPIFromSavedTokenInContext(), cli.NewJsonPrinter(log), args[0], args[1], origin)
		cli.NotifyErrorsWithRetry(err, log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(mapGroupCmd)
	mapGroupCmd.Annotations = make(map[string]string)
	mapGroupCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"

	mapGroupCmd.Flags().StringVarP(&origin, "origin", "", "", "map uaa group to external group for this origin. Defaults to ldap.")
}
