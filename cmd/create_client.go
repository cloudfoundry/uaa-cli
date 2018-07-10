package cmd

import (
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func arrayify(commaSeparatedStr string) []string {
	if commaSeparatedStr == "" {
		return []string{}
	} else {
		return strings.Split(commaSeparatedStr, ",")
	}
}

func CreateClientPreRunValidations(cfg uaa.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) < 1 {
		return MissingArgumentError("client_id")
	}
	return nil
}

func CreateClientCmd(cm *uaa.ClientManager, clone, clientId, clientSecret, displayName, authorizedGrantTypes, authorities, autoapprove, redirectUri, scope string, accessTokenValidity int64, refreshTokenValidity int64) error {
	var toCreate uaa.Client
	var err error
	if clone != "" {
		toCreate, err = cm.Get(clone)
		if err != nil {
			return errors.New(fmt.Sprintf("The client %v could not be found.", clone))
		}

		toCreate.ClientID = clientId
		toCreate.ClientSecret = clientSecret
		if displayName != "" {
			toCreate.DisplayName = displayName
		}
		if authorizedGrantTypes != "" {
			toCreate.AuthorizedGrantTypes = arrayify(authorizedGrantTypes)
		}
		if authorities != "" {
			toCreate.Authorities = arrayify(authorities)
		}
		if autoapprove != "" {
			toCreate.AutoApprove = arrayify(autoapprove)
		}
		if redirectUri != "" {
			toCreate.RedirectURI = arrayify(redirectUri)
		}
		if scope != "" {
			toCreate.Scope = arrayify(scope)
		}
		if refreshTokenValidity != 0 {
			toCreate.RefreshTokenValidity = refreshTokenValidity
		}
		if accessTokenValidity != 0 {
			toCreate.AccessTokenValidity = accessTokenValidity
		}
	} else {
		toCreate = uaa.Client{}
		toCreate.ClientID = clientId
		toCreate.ClientSecret = clientSecret
		toCreate.DisplayName = displayName
		toCreate.AuthorizedGrantTypes = arrayify(authorizedGrantTypes)
		toCreate.Authorities = arrayify(authorities)
		toCreate.AutoApprove = arrayify(autoapprove)
		toCreate.RedirectURI = arrayify(redirectUri)
		toCreate.Scope = arrayify(scope)
		toCreate.AccessTokenValidity = accessTokenValidity
		toCreate.RefreshTokenValidity = refreshTokenValidity
	}

	validationErr := toCreate.Validate()
	if validationErr != nil {
		return validationErr
	}

	created, err := cm.Create(toCreate)
	if err != nil {
		return err
	}

	log.Infof("The client %v has been successfully created.", utils.Emphasize(clientId))
	return cli.NewJsonPrinter(log).Print(created)
}

var createClientCmd = &cobra.Command{
	Use:   "create-client CLIENT_ID -s CLIENT_SECRET --authorized_grant_types GRANT_TYPES",
	Short: "Create an OAuth client registration in the UAA",
	Long:  help.CreateClient(),
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(CreateClientPreRunValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cm := &uaa.ClientManager{GetHttpClient(), cfg}
		err := CreateClientCmd(
			cm,
			clone,
			args[0],
			clientSecret,
			displayName,
			authorizedGrantTypes,
			authorities,
			autoapprove,
			redirectUri,
			scope,
			accessTokenValidity,
			refreshTokenValidity)
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(createClientCmd)
	createClientCmd.Annotations = make(map[string]string)
	createClientCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	createClientCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	createClientCmd.Flags().StringVarP(&authorizedGrantTypes, "authorized_grant_types", "", "", "list of grant types allowed with this client")
	createClientCmd.Flags().StringVarP(&authorities, "authorities", "", "", "scopes requested by client during client_credentials grant")
	createClientCmd.Flags().StringVarP(&autoapprove, "autoapprove", "", "", "Scopes that do not require user approval")
	createClientCmd.Flags().StringVarP(&scope, "scope", "", "", "scopes requested by client during authorization_code, implicit, or password grants")
	createClientCmd.Flags().Int64VarP(&accessTokenValidity, "access_token_validity", "", 0, "the time in seconds before issued access tokens expire")
	createClientCmd.Flags().Int64VarP(&refreshTokenValidity, "refresh_token_validity", "", 0, "the time in seconds before issued refrsh tokens expire")
	createClientCmd.Flags().StringVarP(&displayName, "display_name", "", "", "a friendly human-readable name for this client")
	createClientCmd.Flags().StringVarP(&redirectUri, "redirect_uri", "", "", "callback urls allowed for use in authorization_code and implicit grants")
	createClientCmd.Flags().StringVarP(&clone, "clone", "", "", "client_id of client configuration to clone")
	createClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to create the client")
}
