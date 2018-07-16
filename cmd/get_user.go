package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func GetUserCmd(api *uaa.API, printer cli.Printer, username, origin, attributes string) error {
	user, err := api.GetUserByUsername(username, origin, attributes)
	if err != nil {
		return err
	}

	return printer.Print(user)
}

func GetUserValidations(cfg config.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("The positional argument USERNAME must be specified.")
	}
	return nil
}

var getUserCmd = &cobra.Command{
	Use:   "get-user USERNAME",
	Short: "Look up a user by username",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(GetUserValidations(GetSavedConfig(), args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := GetUserCmd(GetAPIFromSavedTokenInContext(), cli.NewJsonPrinter(log), args[0], origin, attributes)
		NotifyErrorsWithRetry(err, log)
	},
}

func init() {
	RootCmd.AddCommand(getUserCmd)
	getUserCmd.Annotations = make(map[string]string)
	getUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"

	getUserCmd.Flags().StringVarP(&origin, "origin", "o", "", `The identity provider in which to search. Examples: uaa, ldap, etc. `)
	getUserCmd.Flags().StringVarP(&attributes, "attributes", "a", "", `include only these comma-separated user attributes to improve query performance`)
	getUserCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in find the user")
}
