package cli

import (
	"code.cloudfoundry.org/uaa-cli/utils"
	"net/url"
	"fmt"
	"os"
)

type Impersonator interface {
	Start()
	Authorize()
	Done() chan string
}

type ImplicitClientImpersonator struct {
	ClientId           string
	TokenFormat        string
	Scope              string
	UaaBaseUrl         string
	Port               int
	Log                utils.Logger
	AuthCallbackServer CallbackServer
	BrowserLauncher    func(string) error
	done               chan string
}

const CallbackCSS = `<style>
	@import url('https://fonts.googleapis.com/css?family=Source+Sans+Pro');
	html {
		background: #f8f8f8;
		font-family: "Source Sans Pro", sans-serif;
	}
</style>`
const implicitCallbackJS = `<script>
	// This script is needed to send the token fragment from the browser back to
	// the local server. Browsers remove everything after the # before issuing
	// requests so we have to convert these fragments into query params.
	var req = new XMLHttpRequest();
	req.open("GET", "/" + location.hash.replace("#","?"));
	req.send();
</script>`
const implicitCallbackHTML = `<body>
	<h1>Implicit Grant: Success</h1>
	<p>The UAA redirected you to this page with an access token.</p>
	<p> The token has been added to the CLI's active context. You may close this window.</p>
</body>`

func NewImplicitClientImpersonator(clientId,
	uaaBaseUrl string,
	tokenFormat string,
	scope string,
	port int,
	log utils.Logger,
	launcher func(string) error) ImplicitClientImpersonator {

	impersonator := ImplicitClientImpersonator{
		ClientId:        clientId,
		UaaBaseUrl:      uaaBaseUrl,
		TokenFormat:     tokenFormat,
		Scope:           scope,
		Port:            port,
		BrowserLauncher: launcher,
		Log:             log,
		done:            make(chan string),
	}

	callbackServer := NewAuthCallbackServer(implicitCallbackHTML, CallbackCSS, implicitCallbackJS, log, port)
	callbackServer.SetHangupFunc(func(done chan string, values url.Values) {
		token := values.Get("access_token")
		if token != "" {
			done <- token
		}
	})
	impersonator.AuthCallbackServer = callbackServer

	return impersonator
}

func (ici ImplicitClientImpersonator) Start() {
	ici.AuthCallbackServer.Start(ici.Done())
}
func (ici ImplicitClientImpersonator) Authorize() {
	requestValues := url.Values{}
	requestValues.Add("response_type", "token")
	requestValues.Add("client_id", ici.ClientId)
	requestValues.Add("scope", ici.Scope)
	requestValues.Add("token_format", ici.TokenFormat)
	requestValues.Add("redirect_uri", fmt.Sprintf("http://localhost:%v", ici.Port))

	authUrl, err := utils.BuildUrl(ici.UaaBaseUrl, "/oauth/authorize")
	if err != nil {
		ici.Log.Error("Something went wrong while building the authorization URL.")
		os.Exit(1)
	}
	authUrl.RawQuery = requestValues.Encode()

	ici.Log.Info("Launching browser window to " + authUrl.String())
	ici.BrowserLauncher(authUrl.String())
}
func (ici ImplicitClientImpersonator) Done() chan string {
	return ici.done
}
