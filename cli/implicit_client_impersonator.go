package cli

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"code.cloudfoundry.org/uaa-cli/utils"
	"golang.org/x/oauth2"
	"time"
)

type ClientImpersonator interface {
	Start()
	Authorize()
	Done() chan oauth2.Token
}

type ImplicitClientImpersonator struct {
	ClientID           string
	TokenFormat        string
	Scope              string
	UaaBaseURL         string
	Port               int
	Log                Logger
	AuthCallbackServer CallbackServer
	BrowserLauncher    func(string) error
	done               chan oauth2.Token
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
	uaaBaseURL string,
	tokenFormat string,
	scope string,
	port int,
	log Logger,
	launcher func(string) error) ImplicitClientImpersonator {

	impersonator := ImplicitClientImpersonator{
		ClientID:        clientId,
		UaaBaseURL:      uaaBaseURL,
		TokenFormat:     tokenFormat,
		Scope:           scope,
		Port:            port,
		BrowserLauncher: launcher,
		Log:             log,
		done:            make(chan oauth2.Token),
	}

	callbackServer := NewAuthCallbackServer(implicitCallbackHTML, CallbackCSS, implicitCallbackJS, log, port)
	callbackServer.SetHangupFunc(func(done chan url.Values, values url.Values) {
		token := values.Get("access_token")
		if token != "" {
			done <- values
		}
	})
	impersonator.AuthCallbackServer = callbackServer

	return impersonator
}

func (ici ImplicitClientImpersonator) Start() {
	go func() {
		urlValues := make(chan url.Values)
		go ici.AuthCallbackServer.Start(urlValues)
		values := <-urlValues
		response := oauth2.Token{
			AccessToken: values.Get("access_token"),
			TokenType:   values.Get("token_type"),
		}
		response = *response.WithExtra(map[string]interface{}{
			"scope": values.Get("scope"),
			"jti":   values.Get("jti"),
		})

		expiresIn, err := strconv.Atoi(values.Get("expires_in"))
		if err == nil {
			response.Expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
		}
		ici.Done() <- response
	}()
}
func (ici ImplicitClientImpersonator) Authorize() {
	requestValues := url.Values{}
	requestValues.Add("response_type", "token")
	requestValues.Add("client_id", ici.ClientID)
	requestValues.Add("scope", ici.Scope)
	requestValues.Add("token_format", ici.TokenFormat)
	requestValues.Add("redirect_uri", fmt.Sprintf("http://localhost:%v", ici.Port))

	authUrl, err := utils.BuildUrl(ici.UaaBaseURL, "/oauth/authorize")
	if err != nil {
		ici.Log.Error("Something went wrong while building the authorization URL.")
		os.Exit(1)
	}
	authUrl.RawQuery = requestValues.Encode()

	ici.Log.Info("Launching browser window to " + authUrl.String())
	ici.BrowserLauncher(authUrl.String())
}
func (ici ImplicitClientImpersonator) Done() chan oauth2.Token {
	return ici.done
}
