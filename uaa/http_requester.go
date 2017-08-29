package uaa

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"fmt"
)

type Requester interface {
	Get(client *http.Client, context UaaContext, path string, query string) ([]byte, error)
	PostForm(client *http.Client, context UaaContext, path string, query string, body map[string]string) ([]byte, error)
	PostJson(client *http.Client, context UaaContext, path string, query string, body interface{}) ([]byte, error)
}

type UnauthenticatedRequester struct {}
type AuthenticatedRequester struct {}

func is2XX(status int) bool {
	if status >= 200 && status < 300 {
		return true
	}
	return false
}

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

	if !is2XX(resp.StatusCode) {
		return []byte{}, requestError(req.URL.String())
	}
	return bytes, nil
}

func (ug UnauthenticatedRequester) GetBytes(client *http.Client, config Config, path string, query string) ([]byte, error) {
	req, err := UnauthenticatedRequestFactory{}.Get(config.GetActiveTarget(), path, query)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config, path, query)
}

func (ag AuthenticatedRequester) GetBytes(client *http.Client, config Config, path string, query string) ([]byte, error) {
	req, err := AuthenticatedRequestFactory{}.Get(config.GetActiveTarget(), path, query)
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

func (ug UnauthenticatedRequester) PostForm(client *http.Client, config Config, path string, query string, body map[string]string) ([]byte, error) {
	data := mapToUrlValues(body)

	req, err := UnauthenticatedRequestFactory{}.PostForm(config.GetActiveTarget(), path, query, &data)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config, path, query)
}

func (ag AuthenticatedRequester) PostForm(client *http.Client, config Config, path string, query string, body map[string]string) ([]byte, error) {
	data := mapToUrlValues(body)

	req, err := AuthenticatedRequestFactory{}.PostForm(config.GetActiveTarget(), path, query, &data)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config, path, query)
}

func (ug UnauthenticatedRequester) PostJson(client *http.Client, config Config, path string, query string, body interface{}) ([]byte, error) {
	req, err := UnauthenticatedRequestFactory{}.PostJson(config.GetActiveTarget(), path, query, body)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config, path, query)
}

func (ag AuthenticatedRequester) PostJson(client *http.Client, config Config, path string, query string, body interface{}) ([]byte, error) {
	req, err := AuthenticatedRequestFactory{}.PostJson(config.GetActiveTarget(), path, query, body)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config, path, query)
}