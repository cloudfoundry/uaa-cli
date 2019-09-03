package cmd

import (
	"errors"

	"code.cloudfoundry.org/uaa-cli/cli"
	cli_config "code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func RemoveMemberPreRunValidations(config cli_config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(config); err != nil {
		return err
	}

	if len(args) != 2 {
		return errors.New("The positional arguments GROUPNAME and USERNAME must be specified.")
	}

	return nil
}

func RemoveMemberCmd(api *uaa.API, groupName, username string, log cli.Logger) error {
	group, err := api.GetGroupByName(groupName, "")
	if err != nil {
		return err
	}

	user, err := api.GetUserByUsername(username, "", "")
	if err != nil {
		return err
	}

	err = api.RemoveGroupMember(group.ID, user.ID, "", "")
	if err != nil {
		return err
	}

	log.Infof("User %v successfully removed from group %v", utils.Emphasize(username), utils.Emphasize(groupName))

	return nil
}

var removeMemberCmd = &cobra.Command{
	Use:   "remove-member GROUPNAME USERNAME",
	Short: "Remove a user from a group",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(RemoveMemberPreRunValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		userName := args[1]
		cli.NotifyErrorsWithRetry(RemoveMemberCmd(GetAPIFromSavedTokenInContext(), groupName, userName, log), log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(removeMemberCmd)
	removeMemberCmd.Annotations = make(map[string]string)
	removeMemberCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"
}
