package cmd

import (
	"errors"
	"strings"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

var delAttrs []string

func UpdateUserCmd(api *uaa.API, printer cli.Printer, username, familyName, givenName, origin string, emails []string, phones []string, delAttrs []string) error {
	// First, get the existing user
	user, err := api.GetUserByUsername(username, origin, "")
	if err != nil {
		return err
	}

	// Create a copy of the user for updating
	toUpdate := *user

	// Update fields if provided
	if familyName != "" || givenName != "" {
		if toUpdate.Name == nil {
			toUpdate.Name = &uaa.UserName{}
		}
		if familyName != "" {
			toUpdate.Name.FamilyName = familyName
		}
		if givenName != "" {
			toUpdate.Name.GivenName = givenName
		}
	}

	if len(emails) > 0 {
		toUpdate.Emails = buildEmails(emails)
	}

	if len(phones) > 0 {
		toUpdate.PhoneNumbers = buildPhones(phones)
	}

	// Handle attribute deletion
	if len(delAttrs) > 0 {
		for _, attr := range delAttrs {
			switch strings.ToLower(attr) {
			case "phonenumbers", "phone", "phones":
				toUpdate.PhoneNumbers = nil
			case "emails", "email":
				// Don't allow clearing all emails as it may break the user
				log.Infof("Warning: Cannot delete all emails as it may make the user unusable")
			case "name", "familyname", "givenname":
				if strings.ToLower(attr) == "name" || strings.ToLower(attr) == "familyname" {
					if toUpdate.Name != nil {
						toUpdate.Name.FamilyName = ""
					}
				}
				if strings.ToLower(attr) == "name" || strings.ToLower(attr) == "givenname" {
					if toUpdate.Name != nil {
						toUpdate.Name.GivenName = ""
					}
				}
			}
		}
	}

	// Update the user
	updatedUser, err := api.UpdateUser(toUpdate)
	if err != nil {
		return err
	}

	log.Infof("Account for user %v successfully updated.", utils.Emphasize(updatedUser.Username))
	return printer.Print(updatedUser)
}

func UpdateUserValidation(cfg config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("The positional argument USERNAME must be specified.")
	}
	return nil
}

var updateUserCmd = &cobra.Command{
	Use:   "update-user USERNAME",
	Short: "Update a user account",
	PreRun: func(cmd *cobra.Command, args []string) {
		cli.NotifyValidationErrors(UpdateUserValidation(GetSavedConfig(), args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()

		if zoneSubdomain == "" {
			zoneSubdomain = cfg.ZoneSubdomain
		}
		api := GetAPIFromSavedTokenInContext()
		err := UpdateUserCmd(api, cli.NewJsonPrinter(log), args[0], familyName, givenName, origin, emails, phoneNumbers, delAttrs)
		cli.NotifyErrorsWithRetry(err, log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(updateUserCmd)
	updateUserCmd.Annotations = make(map[string]string)
	updateUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"

	updateUserCmd.Flags().StringVarP(&familyName, "family_name", "", "", "family name")
	updateUserCmd.Flags().StringVarP(&givenName, "given_name", "", "", "given name")
	updateUserCmd.Flags().StringVarP(&origin, "origin", "o", "", "user origin")
	updateUserCmd.Flags().StringSliceVarP(&emails, "emails", "", []string{}, "email addresses (multiple may be specified)")
	updateUserCmd.Flags().StringSliceVarP(&phoneNumbers, "phones", "", []string{}, "phone numbers (multiple may be specified)")
	updateUserCmd.Flags().StringSliceVarP(&delAttrs, "del_attrs", "", []string{}, "attributes to remove (phoneNumbers, name, etc.)")
	updateUserCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to update the user")
}