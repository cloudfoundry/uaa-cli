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
	Get(Target, string, string) (*http.Request, error)
	PostForm(Target, string, string, *url.Values) (*http.Request, error)
}

type UnauthenticatedRequestFactory struct {}
type AuthenticatedRequestFactory struct {}

func (urf UnauthenticatedRequestFactory) Get(target Target, path string, query string) (*http.Request, error) {
	targetUrl, err := utils.BuildUrl(target.BaseUrl, path)
	if err != nil {
		return nil, err
	}
	targetUrl.RawQuery = query

	req, err := http.NewRequest("GET", targetUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept","application/json")

	return req, nil
}

func (urf UnauthenticatedRequestFactory) PostForm(target Target, path string, query string, data *url.Values) (*http.Request, error) {
	targetUrl, err := utils.BuildUrl(target.BaseUrl, path)
	if err != nil {
		return nil, err
	}
	targetUrl.RawQuery = query

	bodyBytes := []byte(data.Encode())
	req, err := http.NewRequest("POST", targetUrl.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept","application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(bodyBytes)))

	return req, nil
}

func (arf AuthenticatedRequestFactory) Get(target Target, path string, query string) (*http.Request, error) {
	req, err := UnauthenticatedRequestFactory{}.Get(target, path, query)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "bearer " + target.GetActiveContext().AccessToken)
	if target.GetActiveContext().AccessToken == "" {
		return nil, errors.New("An access token is required to call " + req.URL.String())
	}

	return req, nil
}

func (arf AuthenticatedRequestFactory) PostForm(target Target, path string, query string, data *url.Values) (*http.Request, error) {
	req, err := UnauthenticatedRequestFactory{}.PostForm(target, path, query, data)
	if err != nil {
		return nil, err
	}

	accessToken := target.GetActiveContext().AccessToken
	req.Header.Add("Authorization", "bearer " + accessToken)
	if accessToken == "" {
		return nil, errors.New("An access token is required to call " + req.URL.String())
	}

	return req, nil
}