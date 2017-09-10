package cli

import (
	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type AuthcodeClientImpersonator struct {
	httpClient         *http.Client
	config             uaa.Config
	ClientId           string
	ClientSecret       string
	TokenFormat        string
	Scope              string
	UaaBaseUrl         string
	Port               int
	Log                utils.Logger
	AuthCallbackServer CallbackServer
	BrowserLauncher    func(string) error
	done               chan uaa.TokenResponse
}

const authcodeCallbackHTML = `<body>
	<h1>Authorization Code Grant: Success</h1>
	<p>The UAA redirected you to this page with an authorization code.</p>
	<p>The CLI will exchange this code for an access token. You may close this window.</p>
</body>`

func NewAuthcodeClientImpersonator(
	httpClient *http.Client,
	config uaa.Config,
	clientId,
	clientSecret,
	uaaBaseUrl,
	tokenFormat,
	scope string,
	port int,
	log utils.Logger,
	launcher func(string) error) AuthcodeClientImpersonator {

	impersonator := AuthcodeClientImpersonator{
		httpClient:      httpClient,
		config:          config,
		ClientId:        clientId,
		ClientSecret:    clientSecret,
		UaaBaseUrl:      uaaBaseUrl,
		TokenFormat:     tokenFormat,
		Scope:           scope,
		Port:            port,
		BrowserLauncher: launcher,
		Log:             log,
		done:            make(chan uaa.TokenResponse),
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
		tokenRequester := uaa.AuthorizationCodeClient{ClientId: aci.ClientId, ClientSecret: aci.ClientSecret}
		resp, err := tokenRequester.RequestToken(aci.httpClient, aci.config, uaa.TokenFormat(aci.TokenFormat), code, aci.redirectUri())
		if err != nil {
			fmt.Println(err.Error())
			aci.Log.Error(err.Error())
			os.Exit(1)
		}
		aci.Done() <- resp
	}()
}
func (aci AuthcodeClientImpersonator) Authorize() {
	requestValues := url.Values{}
	requestValues.Add("response_type", "code")
	requestValues.Add("client_id", aci.ClientId)
	requestValues.Add("redirect_uri", aci.redirectUri())

	authUrl, err := utils.BuildUrl(aci.UaaBaseUrl, "/oauth/authorize")
	if err != nil {
		aci.Log.Error("Something went wrong while building the authorization URL.")
		os.Exit(1)
	}
	authUrl.RawQuery = requestValues.Encode()

	aci.Log.Info("Launching browser window to " + authUrl.String())
	aci.BrowserLauncher(authUrl.String())
}
func (aci AuthcodeClientImpersonator) Done() chan uaa.TokenResponse {
	return aci.done
}
func (aci AuthcodeClientImpersonator) redirectUri() string {
	return fmt.Sprintf("http://localhost:%v", aci.Port)
}
