package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func RefreshTokenCmd(cfg config.Config, log cli.Logger, tokenFormat string) error {
	//TODO: use library function to perform conversion
	format := uaa.JSONWebToken
	if tokenFormat == "opaque" {
		format = uaa.OpaqueToken
	}

	api, err := uaa.NewWithRefreshToken(
		cfg.GetActiveTarget().BaseUrl,
		cfg.ZoneSubdomain,
		cfg.GetActiveContext().ClientId,
		clientSecret,
		cfg.GetActiveContext().Token.RefreshToken,
		format,
		cfg.GetActiveTarget().SkipSSLValidation,
	)
	log.Infof("Using the refresh_token from the active context to request a new access token for client %v.", utils.Emphasize(cfg.GetActiveContext().ClientId))
	if err != nil {
		return cli.RequestErrorFromOauthError(err)
	}

	ctx := cfg.GetActiveContext()

	transport := api.AuthenticatedClient.Transport.(*oauth2.Transport)
	token, err := transport.Source.Token()
	if err != nil {
		return cli.RequestErrorFromOauthError(err)
	}

	ctx.Token = *token
	cfg.AddContext(ctx)
	config.WriteConfig(cfg)
	log.Info("Access token successfully fetched and added to active context.")
	return nil
}

func RefreshTokenValidations(cfg config.Config, clientSecret string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if clientSecret == "" {
		return cli.MissingArgumentError("client_secret")
	}
	if cfg.GetActiveContext().ClientId == "" {
		return errors.New("A client_id was not found in the active context.")
	}
	if GetSavedConfig().GetActiveContext().Token.RefreshToken == "" {
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
		cli.NotifyValidationErrors(RefreshTokenValidations(cfg, clientSecret), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyErrorsWithRetry(RefreshTokenCmd(cfg, log, tokenFormat), log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(refreshTokenCmd)
	refreshTokenCmd.Annotations = make(map[string]string)
	refreshTokenCmd.Annotations[TOKEN_CATEGORY] = "true"
	refreshTokenCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	refreshTokenCmd.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
}
