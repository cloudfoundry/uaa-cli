package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func buildEmails(emails []string) []uaa.Email {
	userEmails := []uaa.Email{}
	var newEmail uaa.Email
	for i, email := range emails {
		if i == 0 {
			newEmail = uaa.Email{Primary: utils.NewTrueP(), Value: email}
		} else {
			newEmail = uaa.Email{Primary: utils.NewFalseP(), Value: email}
		}
		userEmails = append(userEmails, newEmail)
	}
	return userEmails
}

func buildPhones(phones []string) []uaa.PhoneNumber {
	userPhoneNumbers := []uaa.PhoneNumber{}
	var phone uaa.PhoneNumber
	for _, number := range phones {
		phone = uaa.PhoneNumber{Value: number}
		userPhoneNumbers = append(userPhoneNumbers, phone)
	}
	return userPhoneNumbers
}

func CreateUserCmd(api *uaa.API, printer cli.Printer, username, familyName, givenName, password, origin string, emails []string, phones []string) error {
	toCreate := uaa.User{
		Username: username,
		Password: password,
		Origin:   origin,
		Name: &uaa.UserName{
			FamilyName: familyName,
			GivenName:  givenName,
		},
	}

	toCreate.Emails = buildEmails(emails)
	toCreate.PhoneNumbers = buildPhones(phones)

	user, err := api.CreateUser(toCreate)
	if err != nil {
		return err
	}

	return printer.Print(user)
}

func CreateUserValidation(cfg config.Config, args []string, familyName, givenName string, emails []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("The positional argument USERNAME must be specified.")
	}
	if familyName == "" {
		return MissingArgumentError("familyName")
	}
	if givenName == "" {
		return MissingArgumentError("givenName")
	}
	if len(emails) == 0 {
		return MissingArgumentError("email")
	}
	return nil
}

var createUserCmd = &cobra.Command{
	Use:     "create-user USERNAME",
	Short:   "Create a user",
	Aliases: []string{"add-user"},
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(CreateUserValidation(GetSavedConfig(), args, familyName, givenName, emails), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		api, err := uaa.NewWithToken(cfg.GetActiveTarget().BaseUrl, cfg.ZoneSubdomain, cfg.GetActiveContext().Token)
		NotifyErrorsWithRetry(err, log)
		api = api.WithSkipSSLValidation(cfg.GetActiveTarget().SkipSSLValidation)
		err = CreateUserCmd(api, cli.NewJsonPrinter(log), args[0], familyName, givenName, userPassword, origin, emails, phoneNumbers)
		NotifyErrorsWithRetry(err, log)
	},
}

func init() {
	RootCmd.AddCommand(createUserCmd)
	createUserCmd.Annotations = make(map[string]string)
	createUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"

	createUserCmd.Flags().StringVarP(&familyName, "familyName", "", "", "family name (required)")
	createUserCmd.Flags().StringVarP(&givenName, "givenName", "", "", "given name (required)")
	createUserCmd.Flags().StringVarP(&userPassword, "password", "p", "", `user password (required for "uaa" origin)`)
	createUserCmd.Flags().StringVarP(&origin, "origin", "o", "uaa", "user origin")
	createUserCmd.Flags().StringSliceVarP(&emails, "email", "", []string{}, "email address (required, multiple may be specified)")
	createUserCmd.Flags().StringSliceVarP(&phoneNumbers, "phone", "", []string{}, "phone number (optional, multiple may be specified)")
	createUserCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to create the user")
}
