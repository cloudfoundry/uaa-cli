package cmd

import (
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func arrayify(commaSeparatedStr string) []string {
	if commaSeparatedStr == "" {
		return []string{}
	} else {
		return strings.Split(commaSeparatedStr, ",")
	}
}

func onlyImplicit(toCreate uaa.UaaClient) bool {
	return len(toCreate.AuthorizedGrantTypes) == 1 && toCreate.AuthorizedGrantTypes[0] == "implicit"
}

var createClientCmd = &cobra.Command{
	Use:   "create-client CLIENT_ID -s CLIENT_SECRET --authorized_grant_types GRANT_TYPES",
	Short: "Create an OAuth client registration in the UAA",
	Long:  help.CreateClient(),
	Run: func(cmd *cobra.Command, args []string) {
		c := GetSavedConfig()
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}

		clientId := args[0]

		var toCreate uaa.UaaClient
		var err error
		if clone != "" {
			toCreate, err = cm.Get(clone)
			if err != nil {
				fmt.Printf("The client %v could not be found.\n", clone)
				TraceRetryMsg(c)
				os.Exit(1)
			}

			toCreate.ClientId = clientId
			toCreate.ClientSecret = clientSecret
			if displayName != "" {
				toCreate.DisplayName = displayName
			}
			if authorizedGrantTypes != "" {
				toCreate.AuthorizedGrantTypes = arrayify(authorizedGrantTypes)
			}
			if authorities != "" {
				toCreate.Authorities = arrayify(authorities)
			}
			if redirectUri != "" {
				toCreate.RedirectUri = arrayify(redirectUri)
			}
			if scope != "" {
				toCreate.Scope = arrayify(scope)
			}
		} else {
			toCreate = uaa.UaaClient{}
			toCreate.ClientId = clientId
			toCreate.ClientSecret = clientSecret
			toCreate.DisplayName = displayName
			toCreate.AuthorizedGrantTypes = arrayify(authorizedGrantTypes)
			toCreate.Authorities = arrayify(authorities)
			toCreate.RedirectUri = arrayify(redirectUri)
			toCreate.Scope = arrayify(scope)
		}

		validationErr := toCreate.PreCreateValidation()
		if validationErr != nil {
			fmt.Println("Error: " + validationErr.Error())
			cmd.Usage()
			os.Exit(1)
		}

		created, err := cm.Create(toCreate)
		if err != nil {
			fmt.Println("An error occurred while creating the client.")
			TraceRetryMsg(c)
			os.Exit(1)
		}

		j, err := json.MarshalIndent(&created, "", "  ")
		if err != nil {
			fmt.Println(err)
			TraceRetryMsg(c)
			os.Exit(1)
		}

		fmt.Printf("The client %v has been successfully created.\n", clientId)
		fmt.Printf("%v\n", string(j))

	},
	Args: func(cmd *cobra.Command, args []string) error {
		EnsureContext()

		if len(args) < 1 {
			return MissingArgument("client_id")
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(createClientCmd)
	createClientCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	createClientCmd.Flags().StringVarP(&authorizedGrantTypes, "authorized_grant_types", "", "", "list of grant types allowed with this client.")
	createClientCmd.Flags().StringVarP(&authorities, "authorities", "", "", "scopes requested by client during client_credentials grant")
	createClientCmd.Flags().StringVarP(&scope, "scope", "", "", "scopes requested by client during authorization_code, implicit, or password grants")
	createClientCmd.Flags().Int32VarP(&accessTokenValidity, "access_token_validity", "", 0, "the time in seconds before issued access tokens expire")
	createClientCmd.Flags().Int32VarP(&refreshTokenValidity, "refresh_token_validity", "", 0, "the time in seconds before issued refrsh tokens expire")
	createClientCmd.Flags().StringVarP(&displayName, "display_name", "", "", "a friendly human-readable name for this client")
	createClientCmd.Flags().StringVarP(&redirectUri, "redirect_uri", "", "", "callback urls allowed for use in authorization_code and implicit grants")
	createClientCmd.Flags().StringVarP(&clone, "clone", "", "", "client_id of client configuration to clone")
	createClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to create the client")
}
