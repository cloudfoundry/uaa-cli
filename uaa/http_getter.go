package uaa

import (
	"io/ioutil"
	"net/http"
	"errors"
	"fmt"
)

type Getter interface {
	GetBytes(client *http.Client, context UaaContext, path string, query string) ([]byte, error)
}


type UnauthenticatedGetter struct {}
type AuthenticatedGetter struct {}

func getAndRead(factory HttpRequestFactory, client *http.Client, config Config, path string, query string) ([]byte, error) {
	req, err := factory.Get(config.Context, path, query)
	if err != nil {
		return []byte{}, err
	}

	if config.Trace {
		logRequest(req)
	}
	resp, err := client.Do(req)
	if (err != nil) {
		return []byte{}, requestError(req.URL.String())
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, unknownError()
	}
	if config.Trace {
		logResponse(resp, bytes)
	}

	if (resp.StatusCode != http.StatusOK) {
		return []byte{}, requestError(req.URL.String())
	}
	return bytes, nil
}

func logResponse(response *http.Response, bytes []byte) {
	fmt.Printf("< %v\n", response.Status)
	logHeaders("<", response.Header)
	fmt.Printf("< %v\n", string(bytes[:]))
	fmt.Println()
}

func logRequest(request *http.Request) {
	fmt.Printf("> %v %v\n", request.Method, request.URL.String())
	logHeaders(">", request.Header)
	fmt.Println()
}

func logHeaders(prefix string, headers map[string][]string) {
	for header, values := range headers {
		for _, value := range values {
			fmt.Printf("%v %v: %v\n", prefix, header, value)
		}
	}
}

func (ug UnauthenticatedGetter) GetBytes(client *http.Client, config Config, path string, query string) ([]byte, error) {
	return getAndRead(UnauthenticatedRequestFactory{}, client, config, path, query)
}

func (ag AuthenticatedGetter) GetBytes(client *http.Client, config Config, path string, query string) ([]byte, error) {
	return getAndRead(AuthenticatedRequestFactory{}, client, config, path, query)
}

func requestError(url string) error {
	return errors.New("An unknown error occurred while calling " + url)
}

func parseError(url string, body []byte) error {
	errorMsg := "An unknown error occurred while parsing response from " + url + ". Response was " + string(body)
	return errors.New(errorMsg)
}

func unknownError() error {
	return errors.New("An unknown error occurred")
}