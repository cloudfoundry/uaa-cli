package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func SaveContext(ctx uaa.UaaContext, log *utils.Logger) {
	c := GetSavedConfig()
	c.AddContext(ctx)
	config.WriteConfig(c)
	log.Info("Access token added to active context.")
}

func addImplicitTokenToContext(clientId string, format string, tokenResponse uaa.TokenResponse, log *utils.Logger) {
	ctx := uaa.UaaContext{
		ClientId:    clientId,
		GrantType:   "implicit",
		AccessToken: tokenResponse.AccessToken,
		TokenType:   uaa.TokenFormat(format),
		Scope:       tokenResponse.Scope,
		JTI:         tokenResponse.JTI,
		ExpiresIn:   tokenResponse.ExpiresIn,
	}

	SaveContext(ctx, log)
}

func ImplicitTokenArgumentValidation(args []string, port int, cmd *cobra.Command) error {
	if len(args) < 1 {
		MissingArgument("client_id", cmd)
	}
	if port == 0 {
		MissingArgument("port", cmd)
	}
	validateTokenFormat(cmd, tokenFormat)
	return nil
}

func ImplicitTokenCommandRun(doneRunning chan bool, clientId string, implicitImp cli.ClientImpersonator, log *utils.Logger) {
	implicitImp.Start()
	implicitImp.Authorize()
	tokenResponse := <-implicitImp.Done()
	addImplicitTokenToContext(clientId, "jwt", tokenResponse, log)
	doneRunning <- true
}

var getImplicitToken = &cobra.Command{
	Use:   "get-implicit-token CLIENT_ID --port REDIRECT_URI_PORT",
	Short: "Obtain a token as an implicit grant client",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureTarget()
	},
	Long: help.ImplicitGrant(),
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan bool)
		baseUrl := GetSavedConfig().GetActiveTarget().BaseUrl
		implicitImp := cli.NewImplicitClientImpersonator(args[0], baseUrl, tokenFormat, scope, port, log, open.Run)
		go ImplicitTokenCommandRun(done, args[0], implicitImp, GetLogger())
		<-done
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return ImplicitTokenArgumentValidation(args, port, cmd)
	},
}

func init() {
	getImplicitToken.Flags().IntVarP(&port, "port", "", 0, "port on which to run local callback server")
	getImplicitToken.Flags().StringVarP(&scope, "scope", "", "openid", "comma-separated scopes to request in token")
	getImplicitToken.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
	getImplicitToken.Annotations = make(map[string]string)
	getImplicitToken.Annotations[TOKEN_CATEGORY] = "true"
	RootCmd.AddCommand(getImplicitToken)
}
