// Copyright © 2017 Jennifer Hamon <jhamon@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
	"github.com/jhamon/uaa-cli/config"
	"github.com/jhamon/uaa-cli/uaa"
	"github.com/jhamon/uaa-cli/help"
)

var (
	password string
	username string
	clientSecret2 string
)

var getPasswordToken = &cobra.Command{
	Use:   "get-password-token CLIENT_ID",
	Short: "obtain a token as a password grant client",
	Long: help.PasswordGrant(),
	Run: func(cmd *cobra.Command, args []string) {
		clientId := args[0]
		requestedType := uaa.OPAQUE

		ccClient := uaa.ResourceOwnerPasswordClient{
			ClientId: clientId,
			ClientSecret: clientSecret2,
			Username: username,
			Password: password,
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
