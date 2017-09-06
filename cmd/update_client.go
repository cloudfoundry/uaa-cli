package cmd

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var updateClientCmd = &cobra.Command{
	Use:   "update-client CLIENT_ID",
	Short: "Update an OAuth client registration in the UAA",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureTarget()
	},
	Run: func(cmd *cobra.Command, args []string) {
		c := GetSavedConfig()
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}

		clientId := args[0]
		toUpdate := uaa.UaaClient{
			ClientId:             clientId,
			DisplayName:          displayName,
			AuthorizedGrantTypes: arrayify(authorizedGrantTypes),
			Authorities:          arrayify(authorities),
			RedirectUri:          arrayify(redirectUri),
			Scope:                arrayify(scope),
		}

		updated, err := cm.Update(toUpdate)
		if err != nil {
			fmt.Println("An error occurred while updating the client.")
			TraceRetryMsg(c)
			os.Exit(1)
		}

		j, err := json.MarshalIndent(&updated, "", "  ")
		if err != nil {
			fmt.Println(err)
			TraceRetryMsg(c)
			os.Exit(1)
		}

		fmt.Printf("The client %v has been successfully updated.", clientId)
		fmt.Printf("\n%v", string(j))

	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return MissingArgument("client_id")
		}
		if clientSecret != "" {
			fmt.Printf(`Client not updated. Please see "uaa set-client-secret -h" to learn more about changing client secrets.`)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(updateClientCmd)
	updateClientCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	updateClientCmd.Flag("client_secret").Hidden = true

	updateClientCmd.Flags().StringVarP(&authorizedGrantTypes, "authorized_grant_types", "", "", "list of grant types allowed with this client.")
	updateClientCmd.Flags().StringVarP(&authorities, "authorities", "", "", "scopes requested by client during client_credentials grant")
	updateClientCmd.Flags().StringVarP(&scope, "scope", "", "", "scopes requested by client during authorization_code, implicit, or password grants")
	updateClientCmd.Flags().Int32VarP(&accessTokenValidity, "access_token_validity", "", 0, "the time in seconds before issued access tokens expire")
	updateClientCmd.Flags().Int32VarP(&refreshTokenValidity, "refresh_token_validity", "", 0, "the time in seconds before issued refrsh tokens expire")
	updateClientCmd.Flags().StringVarP(&displayName, "display_name", "", "", "a friendly human-readable name for this client")
	updateClientCmd.Flags().StringVarP(&redirectUri, "redirect_uri", "", "", "callback urls allowed for use in authorization_code and implicit grants")
	updateClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to update the client")
}