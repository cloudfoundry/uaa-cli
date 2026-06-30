package cmd

import (
	"fmt"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func SetPasswordCmd(api *uaa.API, username, origin, attributes, password, zoneID string) error {
	user, err := api.GetUserByUsername(username, origin, attributes)
	if err != nil {
		return err
	}
	if user.Meta == nil {
		return errors.New("The user did not have expected metadata version.")
	}

	err = setPasswordByID(api, user.ID, password, zoneID)
	if err != nil {
		return err
	}

	log.Infof("Password for user %v successfully set.", utils.Emphasize(user.Username))
	return nil
}

// setPasswordByID makes a PUT request to /Users/{id}/password with {"password": "newpassword"}
func setPasswordByID(api *uaa.API, userID, password, zoneID string) error {
	path := fmt.Sprintf("/Users/%s/password", userID)
	data := fmt.Sprintf(`{"password": "%s"}`, password)

	headers := []string{"Content-Type: application/json"}
	if zoneID != "" {
		headers = append(headers, fmt.Sprintf("X-Identity-Zone-Id: %s", zoneID))
	}

	_, _, status, err := api.Curl(path, "PUT", data, headers)
	if err != nil {
		return err
	}

	if status >= 400 {
		return fmt.Errorf("set password failed with status %d", status)
	}

	return nil
}

func SetPasswordValidations(cfg config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("The positional argument USERNAME must be specified.")
	}

	return nil
}

var setPasswordCmd = &cobra.Command{
	Use:   "set-password USERNAME",
	Short: "Set password for a user (admin)",
	PreRun: func(cmd *cobra.Command, args []string) {
		cli.NotifyValidationErrors(SetPasswordValidations(GetSavedConfig(), args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()

		// Get password from flag or prompt if not provided
		if userPassword == "" {
			secret := cli.InteractiveSecret{Prompt: "New password"}
			var err error
			userPassword, err = secret.Get()
			if err != nil {
				cli.NotifyErrorsWithRetry(err, log, GetSavedConfig())
				return
			}
		}

		if userPassword == "" {
			cli.NotifyErrorsWithRetry(errors.New("Password must be specified with --password flag or entered when prompted."), log, GetSavedConfig())
			return
		}

		if zoneSubdomain == "" {
			zoneSubdomain = cfg.ZoneSubdomain
		}

		token := cfg.GetActiveContext().Token
		api, err := uaa.New(
			cfg.GetActiveTarget().BaseUrl,
			uaa.WithToken(&token),
			uaa.WithZoneID(zoneSubdomain),
			uaa.WithSkipSSLValidation(cfg.GetActiveTarget().SkipSSLValidation),
			uaa.WithVerbosity(verbose),
		)
		if err != nil {
			cli.NotifyErrorsWithRetry(err, log, GetSavedConfig())
			return
		}

		err = SetPasswordCmd(api, args[0], origin, attributes, userPassword, zoneSubdomain)
		cli.NotifyErrorsWithRetry(err, log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(setPasswordCmd)
	setPasswordCmd.Annotations = make(map[string]string)
	setPasswordCmd.Annotations[USER_CRUD_CATEGORY] = "true"

	setPasswordCmd.Flags().StringVarP(&userPassword, "password", "p", "", "new password for the user")
	setPasswordCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to set the password")
	setPasswordCmd.Flags().StringVarP(&origin, "origin", "o", "", "the identity provider in which to search. Examples: uaa, ldap, etc.")
}
