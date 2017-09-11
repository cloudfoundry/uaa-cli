package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
	"code.cloudfoundry.org/uaa-cli/help"
)

var refreshTokenCmd = &cobra.Command{
	Use:   "refresh-token -s CLIENT_SECRET",
	Short: "Obtain an access token using the refresh_token grant type",
	Long: help.RefreshToken(),
	Run: func(cmd *cobra.Command, args []string) {
		c := GetSavedConfig()
		ctx := c.GetActiveContext()
		refreshClient := uaa.RefreshTokenClient{
			ClientId:     ctx.ClientId,
			ClientSecret: clientSecret,
		}
		log.Info("Using the refresh_token and client_id from the active context to request new access token.")
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
			log.Error("A client_id was not found in the active context.")
			cmd.Help()
			os.Exit(1)
		}
		if GetSavedConfig().GetActiveContext().RefreshToken == "" {
			log.Error("A refresh_token was not found in the active context.")
			cmd.Help()
			os.Exit(1)
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
