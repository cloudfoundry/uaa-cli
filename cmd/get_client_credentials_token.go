package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"errors"
	"github.com/spf13/cobra"
	"net/http"
)

func GetClientCredentialsTokenValidations(cfg uaa.Config, args []string, clientSecret string) error {
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

func GetClientCredentialsTokenCmd(cfg uaa.Config, httpClient *http.Client, clientId, clientSecret string) error {
	ccClient := uaa.ClientCredentialsClient{ClientId: clientId, ClientSecret: clientSecret}
	tokenResponse, err := ccClient.RequestToken(httpClient, cfg, uaa.TokenFormat(tokenFormat))
	if err != nil {
		return errors.New("An error occurred while fetching token.")
	}

	activeContext := cfg.GetActiveContext()
	activeContext.GrantType = uaa.CLIENT_CREDENTIALS
	activeContext.ClientId = clientId
	activeContext.TokenResponse = tokenResponse

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
		NotifyErrorsWithRetry(GetClientCredentialsTokenCmd(cfg, GetHttpClient(), args[0], clientSecret), cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(getClientCredentialsTokenCmd)
	getClientCredentialsTokenCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	getClientCredentialsTokenCmd.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
	getClientCredentialsTokenCmd.Annotations = make(map[string]string)
	getClientCredentialsTokenCmd.Annotations[TOKEN_CATEGORY] = "true"
}
