package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jhamon/uaa-cli/help"
	"github.com/jhamon/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	authorizedGrantTypes string
	authorities          string
	accessTokenValidity  int32
	refreshTokenValidity int32
	displayName          string
	autoapprove          string
	scope                string
	redirectUri          string
)

func arrayify(commaSeparatedStr string) []string {
	if commaSeparatedStr == "" {
		return []string{}
	} else {
		return strings.Split(commaSeparatedStr, ",")
	}
}

var createClientCmd = &cobra.Command{
	Use:   "create-client CLIENT_ID",
	Short: "Create an OAuth client registration in the UAA",
	Long:  help.CreateClient(),
	Run: func(cmd *cobra.Command, args []string) {
		c := GetSavedConfig()
		cm := &uaa.ClientManager{GetHttpClient(), GetSavedConfig()}

		clientId := args[0]
		toCreate := uaa.UaaClient{
			ClientId:             clientId,
			ClientSecret:         clientSecret,
			DisplayName:          displayName,
			AuthorizedGrantTypes: arrayify(authorizedGrantTypes),
			Authorities:          arrayify(authorities),
			Autoapprove:          arrayify(autoapprove),
			RedirectUri:          arrayify(redirectUri),
			Scope:                arrayify(scope),
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

		fmt.Printf("The client %v has been successfully created.", clientId)
		fmt.Printf("\n%v", string(j))

	},
	Args: func(cmd *cobra.Command, args []string) error {
		EnsureTarget()

		if len(args) < 1 {
			return MissingArgument("client_id")
		}
		if authorizedGrantTypes == "" {
			return MissingArgument("authorized_grant_types")
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
	createClientCmd.Flags().StringVarP(&autoapprove, "autoapprove", "", "", "list of scopes that do not require user approval")
	createClientCmd.Flags().StringVarP(&redirectUri, "redirect_uri", "", "", "callback urls allowed for use in authorization_code and implicit grants")
}
