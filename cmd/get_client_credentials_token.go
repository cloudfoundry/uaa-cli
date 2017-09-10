package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
)

var getClientCredentialsTokenCmd = &cobra.Command{
	Use:   "get-client-credentials-token CLIENT_ID -s CLIENT_SECRET",
	Short: "Obtain a token using the client_credentials grant type",
	Long:  help.ClientCredentials(),
	Run: func(cmd *cobra.Command, args []string) {
		ccClient := uaa.ClientCredentialsClient{ClientId: args[0], ClientSecret: clientSecret}
		c := GetSavedConfig()
		tokenResponse, err := ccClient.RequestToken(GetHttpClient(), c, uaa.TokenFormat(tokenFormat))
		if err != nil {
			log.Error("An error occurred while fetching token.")
			TraceRetryMsg(c)
			os.Exit(1)
		}

		activeContext := c.GetActiveContext()
		activeContext.GrantType = uaa.CLIENT_CREDENTIALS
		activeContext.ClientId = args[0]
		activeContext.TokenResponse = tokenResponse

		c.AddContext(activeContext)
		config.WriteConfig(c)
		log.Info("Access token successfully fetched and added to context.")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		EnsureTarget()

		if len(args) < 1 {
			MissingArgument("client_id", cmd)
		}
		if clientSecret == "" {
			MissingArgument("client_secret", cmd)
		}
		validateTokenFormat(cmd, tokenFormat)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(getClientCredentialsTokenCmd)
	getClientCredentialsTokenCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	getClientCredentialsTokenCmd.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
	getClientCredentialsTokenCmd.Annotations = make(map[string]string)
	getClientCredentialsTokenCmd.Annotations[TOKEN_CATEGORY] = "true"
}
