package cmd

import (
	"fmt"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
)

var clientSecret string

var getClientCredentialsTokenCmd = &cobra.Command{
	Use:   "get-client-credentials-token CLIENT_ID -s CLIENT_SECRET",
	Short: "obtain a token as a client_credentials grant client",
	Long:  help.ClientCredentials(),
	Run: func(cmd *cobra.Command, args []string) {
		ccClient := uaa.ClientCredentialsClient{ClientId: args[0], ClientSecret: clientSecret}
		c := GetSavedConfig()
		token, err := ccClient.RequestToken(GetHttpClient(), c, uaa.OPAQUE)
		if err != nil {
			fmt.Println("An error occurred while fetching token.")
			TraceRetryMsg(c)
			os.Exit(1)
		}

		activeContext := c.GetActiveContext()
		activeContext.AccessToken = token.AccessToken
		activeContext.ClientId = args[0]
		activeContext.GrantType = uaa.CLIENT_CREDENTIALS
		activeContext.TokenType = uaa.OPAQUE
		activeContext.JTI = token.JTI
		activeContext.ExpiresIn = token.ExpiresIn
		activeContext.Scope = token.Scope
		c.AddContext(activeContext)
		config.WriteConfig(c)
		fmt.Println("Access token successfully fetched and added to context.")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		EnsureTarget()

		if len(args) < 1 {
			return MissingArgument("client_id")
		}
		if clientSecret == "" {
			return MissingArgument("client_secret")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(getClientCredentialsTokenCmd)
	getClientCredentialsTokenCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
}
