package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	cli_config "code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func AddMemberPreRunValidations(config cli_config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(config); err != nil {
		return err
	}

	if len(args) != 2 {
		return errors.New("The positional arguments GROUPNAME and USERNAME must be specified.")
	}

	return nil
}

func AddMemberCmd(api *uaa.API, groupName, username string, log cli.Logger) error {
	group, err := api.GetGroupByName(groupName, "")
	if err != nil {
		return err
	}

	user, err := api.GetUserByUsername(username, "", "")
	if err != nil {
		return err
	}

	err = api.AddGroupMember(group.ID, user.ID, "", "")
	if err != nil {
		return err
	}

	log.Infof("User %v successfully added to group %v", utils.Emphasize(username), utils.Emphasize(groupName))

	return nil
}

var addMemberCmd = &cobra.Command{
	Use:   "add-member GROUPNAME USERNAME",
	Short: "Add a user to a group",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(AddMemberPreRunValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		userName := args[1]
		cli.NotifyErrorsWithRetry(AddMemberCmd(GetAPIFromSavedTokenInContext(), groupName, userName, log), log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(addMemberCmd)
	addMemberCmd.Annotations = make(map[string]string)
	addMemberCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"
}
