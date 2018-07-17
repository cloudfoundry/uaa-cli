package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func DeactivateUserCmd(api *uaa.API, username, origin, attributes string) error {
	user, err := api.GetUserByUsername(username, origin, attributes)
	if err != nil {
		return err
	}
	if user.Meta == nil {
		return errors.New("The user did not have expected metadata version.")
	}
	err = api.DeactivateUser(user.ID, user.Meta.Version)
	if err != nil {
		return err
	}
	log.Infof("Account for user %v successfully deactivated.", utils.Emphasize(user.Username))

	return nil
}

func DeactivateUserValidations(cfg config.Config, args []string) error {
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

		if zoneSubdomain == "" {
			zoneSubdomain = cfg.ZoneSubdomain
		}
		api, err := uaa.NewWithToken(cfg.GetActiveTarget().BaseUrl, zoneSubdomain, cfg.GetActiveContext().Token)
		NotifyErrorsWithRetry(err, log)
		api.SkipSSLValidation = cfg.GetActiveTarget().SkipSSLValidation
		err = DeactivateUserCmd(api, args[0], origin, attributes)
		NotifyErrorsWithRetry(err, log)
	},
}

func init() {
	RootCmd.AddCommand(deactivateUserCmd)
	deactivateUserCmd.Annotations = make(map[string]string)
	deactivateUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"

	deactivateUserCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain from which to deactivate the user")

}
