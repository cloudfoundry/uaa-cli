package cli

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/cloudfoundry-community/go-uaa"
	"golang.org/x/oauth2"
)

type AuthcodeClientImpersonator struct {
	config             config.Config
	ClientID           string
	ClientSecret       string
	TokenFormat        string
	Scope              string
	UaaBaseURL         string
	Port               int
	Log                Logger
	AuthCallbackServer CallbackServer
	BrowserLauncher    func(string) error
	done               chan oauth2.Token
}

const authcodeCallbackHTML = `<body>
	<h1>Authorization Code Grant: Success</h1>
	<p>The UAA redirected you to this page with an authorization code.</p>
	<p>The CLI will exchange this code for an access token. You may close this window.</p>
</body>`

func NewAuthcodeClientImpersonator(
	config config.Config,
	clientId,
	clientSecret,
	tokenFormat,
	scope string,
	port int,
	log Logger,
	launcher func(string) error) AuthcodeClientImpersonator {

	impersonator := AuthcodeClientImpersonator{
		config:          config,
		ClientID:        clientId,
		ClientSecret:    clientSecret,
		TokenFormat:     tokenFormat,
		Scope:           scope,
		Port:            port,
		BrowserLauncher: launcher,
		Log:             log,
		done:            make(chan oauth2.Token),
	}

	callbackServer := NewAuthCallbackServer(authcodeCallbackHTML, CallbackCSS, "", log, port)
	callbackServer.SetHangupFunc(func(done chan url.Values, values url.Values) {
		token := values.Get("code")
		if token != "" {
			done <- values
		}
	})
	impersonator.AuthCallbackServer = callbackServer

	return impersonator
}

func (aci AuthcodeClientImpersonator) Start() {
	go func() {
		urlValues := make(chan url.Values)
		go aci.AuthCallbackServer.Start(urlValues)
		values := <-urlValues
		code := values.Get("code")

		tokenFormat := uaa.JSONWebToken //TODO: Use aci tokenformat to convert from string to int

		redirectUrl, err := url.ParseRequestURI(aci.redirectUri())
		NotifyErrorsWithRetry(err, aci.Log, aci.config)

		api, err := uaa.New(
			aci.config.GetActiveTarget().BaseUrl,
			uaa.WithAuthorizationCode(
				aci.ClientID,
				aci.ClientSecret,
				code,
				tokenFormat,
				redirectUrl,
			),
			uaa.WithZoneID(aci.config.ZoneSubdomain),
			uaa.WithSkipSSLValidation(aci.config.GetActiveTarget().SkipSSLValidation),
			uaa.WithVerbosity(aci.config.Verbose),
		)
		NotifyErrorsWithRetry(err, aci.Log, aci.config)
		token, err := api.Token(context.Background())
		NotifyErrorsWithRetry(err, aci.Log, aci.config)

		aci.Done() <- *token
	}()
}
func (aci AuthcodeClientImpersonator) Authorize() {
	requestValues := url.Values{}
	requestValues.Add("response_type", "code")
	requestValues.Add("client_id", aci.ClientID)
	requestValues.Add("scope", aci.Scope)
	requestValues.Add("redirect_uri", aci.redirectUri())

	authUrl, err := utils.BuildUrl(aci.config.GetActiveTarget().BaseUrl, "/oauth/authorize")
	if err != nil {
		aci.Log.Error("Something went wrong while building the authorization URL.")
		os.Exit(1)
	}
	authUrl.RawQuery = requestValues.Encode()

	aci.Log.Info("Launching browser window to " + authUrl.String() + " where the user should login and grant approvals")
	aci.BrowserLauncher(authUrl.String())
}
func (aci AuthcodeClientImpersonator) Done() chan oauth2.Token {
	return aci.done
}
func (aci AuthcodeClientImpersonator) redirectUri() string {
	return fmt.Sprintf("http://localhost:%v", aci.Port)
}
