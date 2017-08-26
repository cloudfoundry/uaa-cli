package uaa

import (
	"io/ioutil"
	"net/http"
	"errors"
)

type Getter interface {
	GetBytes(client *http.Client, context UaaContext, path string, query string) ([]byte, error)
}

type UnauthenticatedGetter struct {}
type AuthenticatedGetter struct {}

func getAndRead(factory HttpRequestFactory, client *http.Client, context UaaContext, path string, query string) ([]byte, error) {
	httpClient := &http.Client{}
	req, err := factory.Get(context, path, query)
	if err != nil {
		return []byte{}, err
	}

	resp, err := httpClient.Do(req)
	if (err != nil || resp.StatusCode != http.StatusOK) {
		return []byte{}, requestError(req.URL.String())
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, unknownError()
	}

	return bytes, nil
}

func (ug UnauthenticatedGetter) GetBytes(client *http.Client, context UaaContext, path string, query string) ([]byte, error) {
	return getAndRead(UnauthenticatedRequestFactory{}, client, context, path, query)
}

func (ag AuthenticatedGetter) GetBytes(client *http.Client, context UaaContext, path string, query string) ([]byte, error) {
	return getAndRead(AuthenticatedRequestFactory{}, client, context, path, query)
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