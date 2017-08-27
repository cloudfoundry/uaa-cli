package uaa

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"fmt"
)

type Requester interface {
	GetBytes(client *http.Client, context UaaContext, path string, query string) ([]byte, error)
	PostBytes(client *http.Client, context UaaContext, path string, query string, body map[string]string) ([]byte, error)
}

type UnauthenticatedRequester struct {}
type AuthenticatedRequester struct {}

func doAndRead(req *http.Request, client *http.Client, config Config, path string, query string) ([]byte, error) {
	if config.Trace {
		logRequest(req)
	}
	resp, err := client.Do(req)
	if (err != nil) {
		if config.Trace {
			fmt.Printf("%v\n\n", err)
		}

		return []byte{}, requestError(req.URL.String())
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if config.Trace {
			fmt.Printf("%v\n\n", err)
		}

		return []byte{}, unknownError()
	} else if config.Trace {
		logResponse(resp, bytes)
	}

	if (resp.StatusCode != http.StatusOK) {
		return []byte{}, requestError(req.URL.String())
	}
	return bytes, nil
}

func (ug UnauthenticatedRequester) GetBytes(client *http.Client, config Config, path string, query string) ([]byte, error) {
	req, err := UnauthenticatedRequestFactory{}.Get(config.Context, path, query)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config, path, query)
}

func (ag AuthenticatedRequester) GetBytes(client *http.Client, config Config, path string, query string) ([]byte, error) {
	req, err := AuthenticatedRequestFactory{}.Get(config.Context, path, query)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config, path, query)
}

func mapToUrlValues(body map[string]string) url.Values {
	data := url.Values{}
	for key, val := range body {
		data.Add(key, val)
	}
	return data
}

func (ug UnauthenticatedRequester) PostBytes(client *http.Client, config Config, path string, query string, body map[string]string) ([]byte, error) {
	data := mapToUrlValues(body)

	req, err := UnauthenticatedRequestFactory{}.Post(config.Context, path, query, &data)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config, path, query)
}

func (ag AuthenticatedRequester) PostBytes(client *http.Client, config Config, path string, query string, body map[string]string) ([]byte, error) {
	data := mapToUrlValues(body)

	req, err := AuthenticatedRequestFactory{}.Post(config.Context, path, query, &data)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config, path, query)
}