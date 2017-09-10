package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"errors"
	"github.com/spf13/cobra"
	"os"
)

var refreshTokenCmd = &cobra.Command{
	Use:   "refresh-token -s CLIENT_SECRET",
	Short: "Obtain an access token using the refresh_token grant type",
	Run: func(cmd *cobra.Command, args []string) {
		c := GetSavedConfig()
		ctx := c.GetActiveContext()
		refreshClient := uaa.RefreshTokenClient{
			ClientId:     ctx.ClientId,
			ClientSecret: clientSecret,
		}
		tokenResponse, err := refreshClient.RequestToken(GetHttpClient(), c, uaa.TokenFormat(tokenFormat), ctx.RefreshToken)
		if err != nil {
			log.Error(err.Error())
			log.Error("An error occurred while fetching token.")
			TraceRetryMsg(c)
			os.Exit(1)
		}

		ctx.TokenResponse = tokenResponse
		c.AddContext(ctx)
		config.WriteConfig(c)
		log.Info("Access token successfully fetched and added to active context.")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		EnsureContext()
		if clientSecret == "" {
			MissingArgument("client_secret", cmd)
		}

		if GetSavedConfig().GetActiveContext().ClientId == "" {
			return errors.New("A client_id was not found in the active context.")
		}
		if GetSavedConfig().GetActiveContext().RefreshToken == "" {
			return errors.New("A refresh_token was not found in the active context.")
		}

		validateTokenFormat(cmd, tokenFormat)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(refreshTokenCmd)
	refreshTokenCmd.Annotations = make(map[string]string)
	refreshTokenCmd.Annotations[TOKEN_CATEGORY] = "true"
	refreshTokenCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	refreshTokenCmd.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
}
