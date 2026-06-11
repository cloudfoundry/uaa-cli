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

func UnlockUserCmd(api *uaa.API, username, origin, attributes, zoneID string) error {
	user, err := api.GetUserByUsername(username, origin, attributes)
	if err != nil {
		return err
	}
	if user.Meta == nil {
		return errors.New("The user did not have expected metadata version.")
	}
	
	err = unlockUserByID(api, user.ID, zoneID)
	if err != nil {
		return err
	}
	
	log.Infof("Account for user %v successfully unlocked.", utils.Emphasize(user.Username))
	return nil
}

// unlockUserByID makes a PATCH request to /Users/{id}/status with {"locked": false}
func unlockUserByID(api *uaa.API, userID, zoneID string) error {
	path := fmt.Sprintf("/Users/%s/status", userID)
	data := `{"locked": false}`
	
	headers := []string{"Content-Type: application/json"}
	if zoneID != "" {
		headers = append(headers, fmt.Sprintf("X-Identity-Zone-Id: %s", zoneID))
	}
	
	_, _, status, err := api.Curl(path, "PATCH", data, headers)
	if err != nil {
		return err
	}
	
	if status >= 400 {
		return fmt.Errorf("unlock user failed with status %d", status)
	}
	
	return nil
}

func UnlockUserValidations(cfg config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("The positional argument USERNAME must be specified.")
	}
	return nil
}

var unlockUserCmd = &cobra.Command{
	Use:   "unlock-user USERNAME",
	Short: "Unlock a user account by username",
	PreRun: func(cmd *cobra.Command, args []string) {
		cli.NotifyValidationErrors(UnlockUserValidations(GetSavedConfig(), args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()

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
		
		err = UnlockUserCmd(api, args[0], origin, attributes, zoneSubdomain)
		cli.NotifyErrorsWithRetry(err, log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(unlockUserCmd)
	unlockUserCmd.Annotations = make(map[string]string)
	unlockUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"

	unlockUserCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain from which to unlock the user")
	unlockUserCmd.Flags().StringVarP(&origin, "origin", "o", "", "the identity provider in which to search. Examples: uaa, ldap, etc.")
}