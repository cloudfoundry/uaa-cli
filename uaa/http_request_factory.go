package uaa

import (
	"net/http"
	"github.com/jhamon/uaa-cli/utils"
	"errors"
	"net/url"
	"bytes"
	"strconv"
)

type HttpRequestFactory interface {
	Get(UaaContext, string, string) (*http.Request, error)
	Post(UaaContext, string, string, *url.Values) (*http.Request, error)
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

func (urf UnauthenticatedRequestFactory) Post(context UaaContext, path string, query string, data *url.Values) (*http.Request, error) {
	targetUrl, err := utils.BuildUrl(context.BaseUrl, path)
	if err != nil {
		return nil, err
	}
	targetUrl.RawQuery = query
	target := targetUrl.String()

	bodyBytes := []byte(data.Encode())
	req, err := http.NewRequest("POST", target, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept","application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(bodyBytes)))

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

func (arf AuthenticatedRequestFactory) Post(context UaaContext, path string, query string, data *url.Values) (*http.Request, error) {
	req, err := UnauthenticatedRequestFactory{}.Post(context, path, query, data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "bearer " + context.AccessToken)
	if context.AccessToken == "" {
		return nil, errors.New("An access token is required to call " + req.URL.String())
	}

	return req, nil
}