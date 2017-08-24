package uaa

import (
	"io/ioutil"
	"net/http"
	"github.com/jhamon/uaa-cli/utils"
	"errors"
)

type Getter interface {
	Get(context UaaContext, path string, query string) ([]byte, int, error)
}

type UnauthenticatedGetter struct {}
type AuthenticatedGetter struct {}

func (ug UnauthenticatedGetter) Get(context UaaContext, path string, query string) ([]byte, error) {
	targetUrl, err := utils.BuildUrl(context.BaseUrl, path)
	if err != nil {
		return []byte{}, err
	}
	targetUrl.RawQuery = query
	target := targetUrl.String()

	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Add("Accept","application/json")

	resp, err := httpClient.Do(req)
	if (resp.StatusCode != http.StatusOK || err != nil) {
		return []byte{}, requestError(target)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, unknownError()
	}

	return bytes, nil
}

func (ag AuthenticatedGetter) Get(context UaaContext, path string, query string) ([]byte, error) {
	targetUrl, err := utils.BuildUrl(context.BaseUrl, path)
	if err != nil {
		return []byte{}, err
	}
	targetUrl.RawQuery = query
	target := targetUrl.String()

	if context.AccessToken == "" {
		return []byte{}, errors.New("An access token is required to call " + target)
	}

	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "bearer " + context.AccessToken)

	resp, err := httpClient.Do(req)
	if (resp.StatusCode != http.StatusOK || err != nil) {
		return []byte{}, requestError(target)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, unknownError()
	}

	return bytes, nil
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