package uaa

import (
	"io/ioutil"
	"net/http"
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

func (ug UnauthenticatedGetter) GetBytes(client *http.Client, config Config, path string, query string) ([]byte, error) {
	return getAndRead(UnauthenticatedRequestFactory{}, client, config, path, query)
}

func (ag AuthenticatedGetter) GetBytes(client *http.Client, config Config, path string, query string) ([]byte, error) {
	return getAndRead(AuthenticatedRequestFactory{}, client, config, path, query)
}