package cmd

import (
	"fmt"

	"github.com/jhamon/uaa-cli/config"
	"github.com/jhamon/uaa-cli/help"
	"github.com/jhamon/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
)

var (
	password      string
	username      string
	clientSecret2 string
)

var getPasswordToken = &cobra.Command{
	Use:   "get-password-token CLIENT_ID",
	Short: "obtain a token as a password grant client",
	Long:  help.PasswordGrant(),
	Run: func(cmd *cobra.Command, args []string) {
		clientId := args[0]
		requestedType := uaa.OPAQUE

		ccClient := uaa.ResourceOwnerPasswordClient{
			ClientId:     clientId,
			ClientSecret: clientSecret2,
			Username:     username,
			Password:     password,
		}
		c := GetSavedConfig()
		token, err := ccClient.RequestToken(GetHttpClient(), c, requestedType)
		if err != nil {
			fmt.Println("An error occurred while fetching token.")
			TraceRetryMsg(c)
			os.Exit(1)
		}

		activeContext := c.GetActiveContext()
		activeContext.AccessToken = token.AccessToken
		activeContext.ClientId = clientId
		activeContext.Username = username
		activeContext.GrantType = uaa.PASSWORD
		activeContext.TokenType = requestedType
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
		if clientSecret2 == "" {
			return MissingArgument("client_secret")
		}
		if password == "" {
			return MissingArgument("password")
		}
		if username == "" {
			return MissingArgument("username")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(getPasswordToken)
	getPasswordToken.Flags().StringVarP(&clientSecret2, "client_secret", "s", "", "client secret")
	getPasswordToken.Flags().StringVarP(&username, "username", "u", "", "username")
	getPasswordToken.Flags().StringVarP(&password, "password", "p", "", "user password")
}
