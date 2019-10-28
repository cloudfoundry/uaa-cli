package cmd
import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func UnmapGroupCmd(api *uaa.API, printer cli.Printer, externalGroupName, groupName, origin string) error {
	if origin == "" {
		origin = "ldap"
	}

	group, err := api.GetGroupByName(groupName, "")
	if err != nil {
		return err
	}
	err = api.UnmapGroup(group.ID, externalGroupName, origin)
	if err != nil {
		return err
	}

	log.Infof("Successfully unmapped %v from %v for origin %v", utils.Emphasize(groupName), utils.Emphasize(externalGroupName), utils.Emphasize(origin))
	return nil
}

func UnmapGroupValidations(cfg config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) != 2 {
		return errors.New("The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified.")
	}

	return nil
}

var unmapGroupCmd = &cobra.Command{
	Use:   "unmap-group EXTERNAL_GROUPNAME GROUPNAME",
	Short: "Unmaps an external group defined within an external identity provider from a uaa group (scope)",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(UnmapGroupValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := UnmapGroupCmd(GetAPIFromSavedTokenInContext(), cli.NewJsonPrinter(log), args[0], args[1], origin)
		cli.NotifyErrorsWithRetry(err, log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(unmapGroupCmd)
	unmapGroupCmd.Annotations = make(map[string]string)
	unmapGroupCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"
}
