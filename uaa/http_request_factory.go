package uaa

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jhamon/uaa-cli/utils"
	"net/http"
	"net/url"
	"strconv"
)

type HttpRequestFactory interface {
	Get(Target, string, string) (*http.Request, error)
	Delete(Target, string, string) (*http.Request, error)
	PostForm(Target, string, string, *url.Values) (*http.Request, error)
	PostJson(Target, string, string, interface{}) (*http.Request, error)
	PutJson(Target, string, string, interface{}) (*http.Request, error)
}

type UnauthenticatedRequestFactory struct{}
type AuthenticatedRequestFactory struct{}

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
	req.Header.Add("Accept", "application/json")

	return req, nil
}

func (urf UnauthenticatedRequestFactory) Delete(target Target, path string, query string) (*http.Request, error) {
	targetUrl, err := utils.BuildUrl(target.BaseUrl, path)
	if err != nil {
		return nil, err
	}
	targetUrl.RawQuery = query

	req, err := http.NewRequest("DELETE", targetUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

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
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(bodyBytes)))

	return req, nil
}

func (urf UnauthenticatedRequestFactory) PostJson(target Target, path string, query string, objToJsonify interface{}) (*http.Request, error) {
	targetUrl, err := utils.BuildUrl(target.BaseUrl, path)
	if err != nil {
		return nil, err
	}
	targetUrl.RawQuery = query

	objectJson, err := json.Marshal(objToJsonify)
	if err != nil {
		return nil, err
	}

	bodyBytes := []byte(objectJson)
	req, err := http.NewRequest("POST", targetUrl.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(bodyBytes)))

	return req, nil
}

func (urf UnauthenticatedRequestFactory) PutJson(target Target, path string, query string, objToJsonify interface{}) (*http.Request, error) {
	targetUrl, err := utils.BuildUrl(target.BaseUrl, path)
	if err != nil {
		return nil, err
	}
	targetUrl.RawQuery = query

	objectJson, err := json.Marshal(objToJsonify)
	if err != nil {
		return nil, err
	}

	bodyBytes := []byte(objectJson)
	req, err := http.NewRequest("PUT", targetUrl.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(bodyBytes)))

	return req, nil
}

func addAuthorization(req *http.Request, ctx UaaContext) (*http.Request, error) {
	accessToken := ctx.AccessToken
	req.Header.Add("Authorization", "bearer "+accessToken)
	if accessToken == "" {
		return nil, errors.New("An access token is required to call " + req.URL.String())
	}

	return req, nil
}

func (arf AuthenticatedRequestFactory) Get(target Target, path string, query string) (*http.Request, error) {
	req, err := UnauthenticatedRequestFactory{}.Get(target, path, query)

	if err != nil {
		return nil, err
	}

	return addAuthorization(req, target.GetActiveContext())
}

func (arf AuthenticatedRequestFactory) Delete(target Target, path string, query string) (*http.Request, error) {
	req, err := UnauthenticatedRequestFactory{}.Delete(target, path, query)

	if err != nil {
		return nil, err
	}

	return addAuthorization(req, target.GetActiveContext())
}

func (arf AuthenticatedRequestFactory) PostForm(target Target, path string, query string, data *url.Values) (*http.Request, error) {
	req, err := UnauthenticatedRequestFactory{}.PostForm(target, path, query, data)
	if err != nil {
		return nil, err
	}

	return addAuthorization(req, target.GetActiveContext())
}

func (arf AuthenticatedRequestFactory) PostJson(target Target, path string, query string, objToJsonify interface{}) (*http.Request, error) {
	req, err := UnauthenticatedRequestFactory{}.PostJson(target, path, query, objToJsonify)
	if err != nil {
		return nil, err
	}

	return addAuthorization(req, target.GetActiveContext())
}

func (arf AuthenticatedRequestFactory) PutJson(target Target, path string, query string, objToJsonify interface{}) (*http.Request, error) {
	req, err := UnauthenticatedRequestFactory{}.PutJson(target, path, query, objToJsonify)
	if err != nil {
		return nil, err
	}

	return addAuthorization(req, target.GetActiveContext())
}
