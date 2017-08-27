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
	"errors"
	"os"
	"github.com/jhamon/uaa-cli/config"
	"github.com/jhamon/uaa-cli/uaa"
)

var (
	password string
	username string
	clientSecret2 string
)

var getPasswordToken = &cobra.Command{
	Use:   "get-password-token",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		ccClient := uaa.ResourceOwnerPasswordClient{
			ClientId: args[0],
			ClientSecret: clientSecret2,
			Username: username,
			Password: password,
		}
		c := GetSavedConfig()
		token, err := ccClient.RequestToken(GetHttpClient(), c, uaa.OPAQUE)
		if err != nil {
			fmt.Println("An error occurred while fetching token.")
			TraceRetryMsg(c)
			os.Exit(1)
		}

		activeContext := c.GetActiveContext()
		activeContext.AccessToken = token.AccessToken
		c.AddContext(activeContext)
		config.WriteConfig(c)
		fmt.Println("Access token successfully fetched and added to context.")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		EnsureTarget()

		if len(args) < 1 {
			return errors.New("Missing argument `client_id` must be specified.\n")
		}
		if clientSecret2 == "" {
			return errors.New("Missing argument `client_secret` must be specified.\n")
		}
		if password == "" {
			return errors.New("Missing argument `password` must be specified.\n")
		}
		if username == "" {
			return errors.New("Missing argument `username` must be specified.\n")
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
