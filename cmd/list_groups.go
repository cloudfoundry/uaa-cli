package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func ListGroupValidations(cfg uaa.Config) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	return nil
}

func ListGroupsCmd(gm uaa.GroupManager, printer cli.Printer, filter, sortBy, sortOrder, attributes string) error {
	group, err := gm.List(filter, sortBy, attributes, uaa.SortOrder(sortOrder))
	if err != nil {
		return err
	}
	return printer.Print(group)
}

var listGroupsCmd = &cobra.Command{
	Use:     "list-groups",
	Aliases: []string{"groups", "get-groups", "search-groups"},
	Short:   "Search and list groups with SCIM filters",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(ListGroupValidations(GetSavedConfig()), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		gm := uaa.GroupManager{GetHttpClient(), cfg}
		err := ListGroupsCmd(gm, cli.NewJsonPrinter(log), filter, sortBy, sortOrder, attributes)
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(listGroupsCmd)
	listGroupsCmd.Annotations = make(map[string]string)
	listGroupsCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"

	listGroupsCmd.Flags().StringVarP(&filter, "filter", "", "", `a SCIM filter, or query, e.g. 'id eq "a5e3f9fb-65a0-4033-a86c-11f4712e1fed"'`)
	listGroupsCmd.Flags().StringVarP(&sortBy, "sortBy", "b", "", `the attribute to sort results by, e.g. "created" or "displayName"`)
	listGroupsCmd.Flags().StringVarP(&sortOrder, "sortOrder", "o", "", `how results should be ordered. One of: [ascending, descending]`)
	listGroupsCmd.Flags().StringVarP(&attributes, "attributes", "a", "", `include only these comma-separated attributes to improve query performance`)
	listGroupsCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain from which to list the groups")
}
