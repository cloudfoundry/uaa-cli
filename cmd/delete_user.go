package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func DeleteUserCmd(api *uaa.API, username, origin, attributes string) error {
	user, err := api.GetUserByUsername(username, origin, attributes)
	if err != nil {
		return err
	}

	_, err = api.DeleteUser(user.ID)
	if err != nil {
		return err
	}

	log.Infof("Account for user %v successfully deleted.", utils.Emphasize(user.Username))
	return nil
}

func DeleteUserValidations(cfg config.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("The positional argument USERNAME must be specified.")
	}
	return nil
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete-user USERNAME",
	Short: "Delete a user by username",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := DeleteUserValidations(GetSavedConfig(), args)
		NotifyValidationErrors(err, cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := DeleteUserCmd(GetAPIFromSavedTokenInContext(), args[0], origin, attributes)
		NotifyErrorsWithRetry(err, log)
	},
}

func init() {
	RootCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Annotations = make(map[string]string)
	deleteUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"

	deleteUserCmd.Flags().StringVarP(&origin, "origin", "o", "uaa", "user origin")
}
