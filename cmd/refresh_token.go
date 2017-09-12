package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"net/http"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
)

func RefreshTokenCmd(cfg uaa.Config, httpClient *http.Client, log utils.Logger, tokenFormat string) error {
	ctx := cfg.GetActiveContext()
	refreshClient := uaa.RefreshTokenClient{
		ClientId:     ctx.ClientId,
		ClientSecret: clientSecret,
	}
	log.Infof("Using the refresh_token from the active context to request a new access token for client %v.", ctx.ClientId)
	tokenResponse, err := refreshClient.RequestToken(httpClient, cfg, uaa.TokenFormat(tokenFormat), ctx.RefreshToken)
	if err != nil {
		return err
	}

	ctx.TokenResponse = tokenResponse
	cfg.AddContext(ctx)
	config.WriteConfig(cfg)
	log.Info("Access token successfully fetched and added to active context.")
	return nil
}

func RefreshTokenValidations(cfg uaa.Config, clientSecret string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if clientSecret == "" {
		return MissingArgumentError("client_secret")
	}
	if cfg.GetActiveContext().ClientId == "" {
		return errors.New("A client_id was not found in the active context.")
	}
	if GetSavedConfig().GetActiveContext().RefreshToken == "" {
		return errors.New("A refresh_token was not found in the active context.")
	}

	return validateTokenFormatError(tokenFormat)
}

var refreshTokenCmd = &cobra.Command{
	Use:   "refresh-token -s CLIENT_SECRET",
	Short: "Obtain an access token using the refresh_token grant type",
	Long:  help.RefreshToken(),
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(RefreshTokenValidations(cfg, clientSecret), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyErrorsWithRetry(RefreshTokenCmd(cfg, GetHttpClient(), log, tokenFormat), cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(refreshTokenCmd)
	refreshTokenCmd.Annotations = make(map[string]string)
	refreshTokenCmd.Annotations[TOKEN_CATEGORY] = "true"
	refreshTokenCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	refreshTokenCmd.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
}
