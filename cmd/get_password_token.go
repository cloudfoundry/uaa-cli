package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"errors"
	"github.com/spf13/cobra"
	"net/http"
)

func GetPasswordTokenValidations(cfg uaa.Config, args []string, clientSecret, username, password string) error {
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

func GetPasswordTokenCmd(cfg uaa.Config, httpClient *http.Client, clientId, clientSecret, username, password, tokenFormat string) error {
	requestedType := uaa.TokenFormat(tokenFormat)

	ccClient := uaa.ResourceOwnerPasswordClient{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
	}
	tokenResponse, err := ccClient.RequestToken(httpClient, cfg, requestedType)
	if err != nil {
		return errors.New("An error occurred while fetching token.")
	}

	activeContext := cfg.GetActiveContext()
	activeContext.ClientId = clientId
	activeContext.GrantType = uaa.PASSWORD
	activeContext.Username = username
	activeContext.TokenResponse = tokenResponse
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
		NotifyValidationErrors(GetPasswordTokenValidations(cfg, args, clientSecret, username, password), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyErrorsWithRetry(GetPasswordTokenCmd(cfg, GetHttpClient(), args[0], clientSecret, username, password, tokenFormat), cfg, log)
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
