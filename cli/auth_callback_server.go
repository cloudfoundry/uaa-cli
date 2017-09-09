package cli

import (
	"code.cloudfoundry.org/uaa-cli/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type AuthCallbackServer struct {
	Html       string
	CSS        string
	JavaScript string
	Port       int
	Log        utils.Logger
	Hangup     func(chan string, url.Values)
	srv        *http.Server
}

func NewAuthCallbackServer(html, css, js string, port int) AuthCallbackServer {
	acs := AuthCallbackServer{Html: html, CSS: css, JavaScript: js, Port: port}
	acs.Log = utils.NewLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	acs.Hangup = func(done chan string, vals url.Values) {}
	return acs
}

func (ls *AuthCallbackServer) Start(done chan string) {
	callbackValue := make(chan string)
	serveMux := http.NewServeMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", ls.Port),
		Handler: serveMux,
	}

	go func() {
		value := <-callbackValue
		close(callbackValue)
		srv.Close()
		done <- value
	}()

	attemptHangup := func(queryParams url.Values) {
		time.Sleep(10 * time.Millisecond)
		ls.Hangup(callbackValue, queryParams)
	}

	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, ls.CSS)
		io.WriteString(w, ls.Html)
		io.WriteString(w, ls.JavaScript)
		ls.Log.Infof("Local server received request to %v %v", r.Method, r.RequestURI)

		// This is a goroutine because we want this handleFunc to complete before
		// Server.Close is invoked by listeners on the callbackValue channel.
		go attemptHangup(r.URL.Query())
	})

	ls.Log.Infof("Starting local HTTP server on port %v", ls.Port)
	ls.Log.Info("Waiting for authorization redirect with code from UAA...")
	if err := srv.ListenAndServe(); err != nil {
		ls.Log.Infof("Stopping local HTTP server on port %v", ls.Port)
	}
}
