package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func ListGroupMappingsCmd(api *uaa.API, printer cli.Printer) error {
	group, err := api.ListAllGroupMappings("")
	if err != nil {
		return err
	}
	return printer.Print(group)
}

var listGroupMappingsCmd = &cobra.Command{
	Use:     "list-group-mappings",
	Aliases: []string{},
	Short:   "List all the mappings between uaa scopes and external groups",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(ListGroupValidations(GetSavedConfig()), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := ListGroupMappingsCmd(GetAPIFromSavedTokenInContext(), cli.NewJsonPrinter(log))
		NotifyErrorsWithRetry(err, log)
	},
}

func init() {
	RootCmd.AddCommand(listGroupMappingsCmd)
	listGroupMappingsCmd.Annotations = make(map[string]string)
	listGroupMappingsCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"
}
