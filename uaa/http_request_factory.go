package uaa

import (
	"net/http"
	"github.com/jhamon/uaa-cli/utils"
	"errors"
)

type HttpRequestFactory interface {
	Get(UaaContext, string, string) (*http.Request, error)
}

type UnauthenticatedRequestFactory struct {}
type AuthenticatedRequestFactory struct {}

func (urf UnauthenticatedRequestFactory) Get(context UaaContext, path string, query string) (*http.Request, error) {
	targetUrl, err := utils.BuildUrl(context.BaseUrl, path)
	if err != nil {
		return nil, err
	}
	targetUrl.RawQuery = query
	target := targetUrl.String()

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept","application/json")

	return req, nil
}

func (arf AuthenticatedRequestFactory) Get(context UaaContext, path string, query string) (*http.Request, error) {
	req, err := UnauthenticatedRequestFactory{}.Get(context, path, query)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "bearer " + context.AccessToken)
	if context.AccessToken == "" {
		return nil, errors.New("An access token is required to call " + req.URL.String())
	}

	return req, nil
}