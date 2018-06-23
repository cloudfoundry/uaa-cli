package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"github.com/cloudfoundry-community/go-uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/spf13/cobra"
)

func DeactivateUserCmd(um uaa.UserManager, printer cli.Printer, username, origin, attributes string) error {
	user, err := um.GetByUsername(username, origin, attributes)
	if err != nil {
		return err
	}
	if user.Meta == nil {
		return errors.New("The user did not have expected metadata version.")
	}
	err = um.Deactivate(user.ID, user.Meta.Version)
	if err != nil {
		return err
	}
	log.Infof("Account for user %v successfully deactivated.", utils.Emphasize(user.Username))

	return nil
}

func DeactivateUserValidations(cfg uaa.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("The positional argument USERNAME must be specified.")
	}
	return nil
}

var deactivateUserCmd = &cobra.Command{
	Use:   "deactivate-user USERNAME",
	Short: "Deactivate a user by username",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(DeactivateUserValidations(GetSavedConfig(), args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		um := uaa.UserManager{GetHttpClient(), cfg}
		err := DeactivateUserCmd(um, cli.NewJsonPrinter(log), args[0], origin, attributes)
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(deactivateUserCmd)
	deactivateUserCmd.Annotations = make(map[string]string)
	deactivateUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"

	deactivateUserCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain from which to deactivate the user")

}
