package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"github.com/pkg/errors"
)

func GetClientCredentialsTokenValidations(cfg config.Config, args []string, clientSecret string) error {
	if err := EnsureTargetInConfig(cfg); err != nil {
		return err
	}
	if len(args) < 1 {
		return MissingArgumentError("client_id")
	}
	if clientSecret == "" {
		return MissingArgumentError("client_secret")
	}
	return validateTokenFormatError(tokenFormat)
}

func GetClientCredentialsTokenCmd(cfg config.Config, clientId, clientSecret string) error {
	var uaaTokenFormat uaa.TokenFormat
	//TODO: place in a method to determine uaaTokenFormat
	if uaa.JSONWebToken.String() == tokenFormat {
		uaaTokenFormat = uaa.JSONWebToken
	} else {
		uaaTokenFormat = uaa.OpaqueToken
	}

	api, err := uaa.NewWithClientCredentials(cfg.GetActiveTarget().BaseUrl, cfg.ZoneSubdomain, clientId, clientSecret, uaaTokenFormat)
	if err != nil {
		return errors.Wrap(err, "An error occurred while building API with client credentials.")
	}

	transport := api.AuthenticatedClient.Transport.(*oauth2.Transport)
	token, err := transport.Source.Token()
	if err != nil {
		return errors.Wrap(err, "An error occurred while fetching token.")
	}

	activeContext := cfg.GetActiveContext()
	activeContext.GrantType = config.CLIENT_CREDENTIALS
	activeContext.ClientId = clientId
	activeContext.Token = *token

	cfg.AddContext(activeContext)
	config.WriteConfig(cfg)
	log.Info("Access token successfully fetched and added to context.")
	return nil
}

var getClientCredentialsTokenCmd = &cobra.Command{
	Use:   "get-client-credentials-token CLIENT_ID -s CLIENT_SECRET",
	Short: "Obtain an access token using the client_credentials grant type",
	Long:  help.ClientCredentials(),
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(GetClientCredentialsTokenValidations(cfg, args, clientSecret), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyErrorsWithRetry(GetClientCredentialsTokenCmd(cfg, args[0], clientSecret), log)
	},
}

func init() {
	RootCmd.AddCommand(getClientCredentialsTokenCmd)
	getClientCredentialsTokenCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	getClientCredentialsTokenCmd.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
	getClientCredentialsTokenCmd.Annotations = make(map[string]string)
	getClientCredentialsTokenCmd.Annotations[TOKEN_CATEGORY] = "true"
}
