package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func ListUserValidations(cfg uaa.Config) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	return nil
}

func ListUsersCmd(um uaa.UserManager, printer cli.Printer, filter, sortBy, sortOrder, attributes string) error {
	user, err := um.List(filter, sortBy, attributes, uaa.SortOrder(sortOrder))
	if err != nil {
		return err
	}
	return printer.Print(user)
}

var listUsersCmd = &cobra.Command{
	Use:     "list-users",
	Aliases: []string{"users", "get-users", "search-users"},
	Short:   "Search and list users with SCIM filters",
	Long:    help.ListUsers(),
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(ListUserValidations(GetSavedConfig()), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		um := uaa.UserManager{GetHttpClient(), cfg}
		err := ListUsersCmd(um, cli.NewJsonPrinter(log), filter, sortBy, sortOrder, attributes)
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(listUsersCmd)
	listUsersCmd.Annotations = make(map[string]string)
	listUsersCmd.Annotations[USER_CRUD_CATEGORY] = "true"

	listUsersCmd.Flags().StringVarP(&filter, "filter", "", "", `a SCIM filter, or query, e.g. 'userName eq "bob@example.com"'`)
	listUsersCmd.Flags().StringVarP(&sortBy, "sortBy", "b", "", `the attribute to sort results by, e.g. "created" or "userName"`)
	listUsersCmd.Flags().StringVarP(&sortOrder, "sortOrder", "o", "", `how results should be ordered. One of: [ascending, descending]`)
	listUsersCmd.Flags().StringVarP(&attributes, "attributes", "a", "", `include only these comma-separated user attributes to improve query performance`)
	listUsersCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to list the users")
}
