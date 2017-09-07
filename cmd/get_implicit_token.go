package cmd

import (
	"fmt"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"net/url"
	"os"
)

func addImplicitTokenToContext(clientId string, requestParams url.Values, responseParams url.Values) {
	ctx := uaa.UaaContext{
		ClientId:    clientId,
		GrantType:   "implicit",
		AccessToken: responseParams.Get("access_token"),
		TokenType: uaa.TokenFormat(requestParams.Get("token_format")),
		Scope: requestParams.Get("scope"),
	}
	c := GetSavedConfig()
	c.AddContext(ctx)
	config.WriteConfig(c)
	log.Info("Access token added to active context.")
}

func implicitCallbackJS(port int) string {
	return fmt.Sprintf(`<script>
	// This script is needed to send the token fragment from the browser back to
	// the local server. Browsers remove everything after the # before issuing
	// requests so we have to convert these fragments into query params.
	var req = new XMLHttpRequest();
	req.open("GET", "http://localhost:%v/" + location.hash.replace("#","?"));
	req.send();
</script>`, port)
}

const ImplicitCallbackCSS = `<style>
	@import url('https://fonts.googleapis.com/css?family=Source+Sans+Pro');
	html {
		background: #f8f8f8;
		font-family: "Source Sans Pro", sans-serif;
	}
</style>`
const ImplicitCallbackHTML = `<body>
	<h1>Implicit Grant: Success</h1>
	<p>The UAA redirected you to this page with an access token.</p>
	<p> The token has been added to the CLI's active context. You may close this window.</p>
</body>`

func startHttpServer(port int, done chan url.Values) {
	serveMux := http.NewServeMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: serveMux,
	}

	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, implicitCallbackJS(port))
		io.WriteString(w, ImplicitCallbackCSS)
		io.WriteString(w, ImplicitCallbackHTML)
		log.Infof("Local server received request to %v %v", r.Method, r.RequestURI)
		token := r.URL.Query().Get("access_token")
		if token != "" {
			done <- r.URL.Query()
		}
	})

	log.Infof("Starting local HTTP server on port %v", port)
	log.Info("Waiting for authorization redirect from UAA...")
	if err := srv.ListenAndServe(); err != nil {
		log.Infof("Stopping local HTTP server on port %v", port)
	}
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

func ImplicitTokenCommandRun(doneRunning chan bool, launcher func(string) error, scope string, clientId string, port int) {
	done := make(chan url.Values)
	go startHttpServer(port, done)

	requestValues := url.Values{}
	requestValues.Add("response_type", "token")
	requestValues.Add("client_id", clientId)
	requestValues.Add("scope", scope)
	requestValues.Add("token_format", tokenFormat)
	requestValues.Add("redirect_uri", fmt.Sprintf("http://localhost:%v", port))

	authUrl, err := utils.BuildUrl(GetSavedConfig().GetActiveTarget().BaseUrl, "/oauth/authorize")
	if err != nil {
		log.Error("Something went wrong while building the authorization URL.")
		os.Exit(1)
	}
	authUrl.RawQuery = requestValues.Encode()

	log.Info("Launching browser window to " + authUrl.String())
	launcher(authUrl.String())
	responseValues := <-done
	addImplicitTokenToContext(clientId, requestValues, responseValues)
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
		doneRunning := make(chan bool)
		go ImplicitTokenCommandRun(doneRunning, open.Run, scope, args[0], port)
		<-doneRunning
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return ImplicitTokenArgumentValidation(args, port, cmd)
	},
}

func init() {
	getImplicitToken.Flags().IntVarP(&port, "port", "", 0, "port on which to run local callback server")
	getImplicitToken.Flags().StringVarP(&scope, "scope", "", "openid", "comma-separated scopes to request in token")
	getImplicitToken.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include " + availableFormatsStr())
	RootCmd.AddCommand(getImplicitToken)
}
