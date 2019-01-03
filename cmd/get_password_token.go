package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"errors"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func GetPasswordTokenValidations(cfg config.Config, args []string, username, password string) error {
	if err := EnsureTargetInConfig(cfg); err != nil {
		return err
	}
	if len(args) < 1 {
		return MissingArgumentError("client_id")
	}
	if password == "" {
		return MissingArgumentError("password")
	}
	if username == "" {
		return MissingArgumentError("username")
	}
	return validateTokenFormatError(tokenFormat)
}

func GetPasswordTokenCmd(cfg config.Config, clientId, clientSecret, username, password, tokenFormat string) error {
	requestedType := uaa.OpaqueToken
	if tokenFormat == uaa.JSONWebToken.String() {
		requestedType = uaa.JSONWebToken
	}

	api, err := uaa.NewWithPasswordCredentials(cfg.GetActiveTarget().BaseUrl, cfg.ZoneSubdomain, clientId, clientSecret, username, password, requestedType, cfg.GetActiveTarget().SkipSSLValidation)
	if err != nil {
		return errors.New("An error occurred while fetching token.")
	}
	api.SkipSSLValidation = cfg.GetActiveTarget().SkipSSLValidation

	transport := api.AuthenticatedClient.Transport.(*oauth2.Transport)
	token, err := transport.Source.Token()

	if err != nil {
		log.Info("Unable to retrieve token")
		return errors.New("An error occurred while fetching token.")
	}

	activeContext := cfg.GetActiveContext()
	activeContext.ClientId = clientId
	activeContext.GrantType = config.PASSWORD
	activeContext.Username = username

	activeContext.Token = *token
	cfg.AddContext(activeContext)
	config.WriteConfig(cfg)
	log.Info("Access token successfully fetched and added to context.")
	return nil
}

var getPasswordToken = &cobra.Command{
	Use:   "get-password-token CLIENT_ID -s CLIENT_SECRET -u USERNAME -p PASSWORD",
	Short: "Obtain an access token using the password grant type",
	Long:  help.PasswordGrant(),
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(GetPasswordTokenValidations(cfg, args, username, password), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyErrorsWithRetry(GetPasswordTokenCmd(cfg, args[0], clientSecret, username, password, tokenFormat), log)
	},
}

func init() {
	RootCmd.AddCommand(getPasswordToken)
	getPasswordToken.Annotations = make(map[string]string)
	getPasswordToken.Annotations[TOKEN_CATEGORY] = "true"
	getPasswordToken.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	getPasswordToken.Flags().StringVarP(&username, "username", "u", "", "username")
	getPasswordToken.Flags().StringVarP(&password, "password", "p", "", "user password")
	getPasswordToken.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
}
