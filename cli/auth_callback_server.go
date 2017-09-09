package cli

import (
	"code.cloudfoundry.org/uaa-cli/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type CallbackServer interface {
	Html() string
	CSS() string
	Javascript() string
	Port() int
	Log() utils.Logger
	Hangup(chan string, url.Values)
	Start(chan string)
}

type AuthCallbackServer struct {
	html       string
	css        string
	javascript string
	port       int
	log        utils.Logger
	hangupFunc func(chan string, url.Values)
}

func NewAuthCallbackServer(html, css, js string, log utils.Logger, port int) AuthCallbackServer {
	acs := AuthCallbackServer{html: html, css: css, javascript: js, log: log, port: port}
	acs.SetHangupFunc(func(done chan string, vals url.Values) {})
	return acs
}

func (acs AuthCallbackServer) Html() string {
	return acs.html
}
func (acs AuthCallbackServer) CSS() string {
	return acs.css
}
func (acs AuthCallbackServer) Javascript() string {
	return acs.javascript
}
func (acs AuthCallbackServer) Port() int {
	return acs.port
}
func (acs AuthCallbackServer) Log() utils.Logger {
	return acs.log
}
func (acs AuthCallbackServer) Hangup(done chan string, values url.Values) {
	acs.hangupFunc(done, values)
}
func (acs *AuthCallbackServer) SetHangupFunc(hangupFunc func(chan string, url.Values)) {
	acs.hangupFunc = hangupFunc
}

func (acs AuthCallbackServer) Start(done chan string) {
	callbackValue := make(chan string)
	serveMux := http.NewServeMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", acs.port),
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
		acs.Hangup(callbackValue, queryParams)
	}

	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, acs.css)
		io.WriteString(w, acs.html)
		io.WriteString(w, acs.javascript)
		acs.log.Infof("Local server received request to %v %v", r.Method, r.RequestURI)

		// This is a goroutine because we want this handleFunc to complete before
		// Server.Close is invoked by listeners on the callbackValue channel.
		go attemptHangup(r.URL.Query())
	})

	acs.log.Infof("Starting local HTTP server on port %v", acs.port)
	acs.log.Info("Waiting for authorization redirect with code from UAA...")
	if err := srv.ListenAndServe(); err != nil {
		acs.log.Infof("Stopping local HTTP server on port %v", acs.port)
	}
}

type FakeCallbackServer struct {
	html       string
	css        string
	javascript string
	port       int
	log        utils.Logger
	hangupFunc func(chan string, url.Values)
}

func (fcs FakeCallbackServer) Html() string {
	return fcs.html
}
func (fcs FakeCallbackServer) CSS() string {
	return fcs.css
}
func (fcs FakeCallbackServer) Javascript() string {
	return fcs.javascript
}
func (fcs FakeCallbackServer) Port() int {
	return fcs.port
}
func (fcs FakeCallbackServer) Log() utils.Logger {
	return fcs.log
}
func (fcs FakeCallbackServer) Hangup(done chan string, values url.Values) {
	fcs.hangupFunc(done, values)
}
func (fcs *FakeCallbackServer) SetHangupFunc(hangupFunc func(chan string, url.Values)) {
	fcs.hangupFunc = hangupFunc
}
func (fcs FakeCallbackServer) Start(done chan string) {
	done <- "server was started"
}
