package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func ListMFAProvidersValidations(cfg config.Config) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	return nil
}

func ListMFAProvidersCmd(api *uaa.API, printer cli.Printer) error {
	providers, err := api.ListMFAProviders()
	if err != nil {
		return err
	}
	return printer.Print(providers)
}

var listMFAProvidersCmd = &cobra.Command{
	Use:     "list-mfa-providers",
	Aliases: []string{"users", "get-mfa-providers", "search-mfa-providers"},
	Short:   "Search and list MFA providers",
	Long:    help.ListMFAProviders(),
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(ListMFAProvidersValidations(GetSavedConfig()), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := ListMFAProvidersCmd(GetAPIFromSavedTokenInContext(), cli.NewJsonPrinter(log))
		NotifyErrorsWithRetry(err, log)
	},
}

func init() {
	RootCmd.AddCommand(listMFAProvidersCmd)
	listMFAProvidersCmd.Annotations = make(map[string]string)
	listMFAProvidersCmd.Annotations[USER_CRUD_CATEGORY] = "true"
}
