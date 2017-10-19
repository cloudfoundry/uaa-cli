package uaa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Requester interface {
	Get(client *http.Client, config Config, path string, query string) ([]byte, error)
	Delete(client *http.Client, config Config, path string, query string) ([]byte, error)
	PostForm(client *http.Client, config Config, path string, query string, body map[string]string) ([]byte, error)
	PostJson(client *http.Client, config Config, path string, query string, body interface{}) ([]byte, error)
	PutJson(client *http.Client, config Config, path string, query string, body interface{}) ([]byte, error)
}

type UnauthenticatedRequester struct{}
type AuthenticatedRequester struct{}

func is2XX(status int) bool {
	if status >= 200 && status < 300 {
		return true
	}
	return false
}

func addZoneSwitchHeader(req *http.Request, config *Config) {
	req.Header.Add("X-Identity-Zone-Subdomain", config.ZoneSubdomain)
}

func mapToUrlValues(body map[string]string) url.Values {
	data := url.Values{}
	for key, val := range body {
		data.Add(key, val)
	}
	return data
}

func doAndRead(req *http.Request, client *http.Client, config Config) ([]byte, error) {
	if config.Verbose {
		logRequest(req)
	}

	resp, err := client.Do(req)
	if err != nil {
		if config.Verbose {
			fmt.Printf("%v\n\n", err)
		}

		return []byte{}, requestError(req.URL.String())
	}

	if config.Verbose {
		logResponse(resp)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if config.Verbose {
			fmt.Printf("%v\n\n", err)
		}

		return []byte{}, unknownError()
	}

	if !is2XX(resp.StatusCode) {
		return []byte{}, requestError(req.URL.String())
	}
	return bytes, nil
}

func (ug UnauthenticatedRequester) Get(client *http.Client, config Config, path string, query string) ([]byte, error) {
	req, err := UnauthenticatedRequestFactory{}.Get(config.GetActiveTarget(), path, query)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ag AuthenticatedRequester) Get(client *http.Client, config Config, path string, query string) ([]byte, error) {
	req, err := AuthenticatedRequestFactory{}.Get(config.GetActiveTarget(), path, query)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ug UnauthenticatedRequester) Delete(client *http.Client, config Config, path string, query string) ([]byte, error) {
	req, err := UnauthenticatedRequestFactory{}.Delete(config.GetActiveTarget(), path, query)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ug AuthenticatedRequester) Delete(client *http.Client, config Config, path string, query string) ([]byte, error) {
	req, err := AuthenticatedRequestFactory{}.Delete(config.GetActiveTarget(), path, query)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ug UnauthenticatedRequester) PostForm(client *http.Client, config Config, path string, query string, body map[string]string) ([]byte, error) {
	data := mapToUrlValues(body)

	req, err := UnauthenticatedRequestFactory{}.PostForm(config.GetActiveTarget(), path, query, &data)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ag AuthenticatedRequester) PostForm(client *http.Client, config Config, path string, query string, body map[string]string) ([]byte, error) {
	data := mapToUrlValues(body)

	req, err := AuthenticatedRequestFactory{}.PostForm(config.GetActiveTarget(), path, query, &data)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ug UnauthenticatedRequester) PostJson(client *http.Client, config Config, path string, query string, body interface{}) ([]byte, error) {
	req, err := UnauthenticatedRequestFactory{}.PostJson(config.GetActiveTarget(), path, query, body)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ag AuthenticatedRequester) PostJson(client *http.Client, config Config, path string, query string, body interface{}) ([]byte, error) {
	req, err := AuthenticatedRequestFactory{}.PostJson(config.GetActiveTarget(), path, query, body)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ug UnauthenticatedRequester) PutJson(client *http.Client, config Config, path string, query string, body interface{}) ([]byte, error) {
	req, err := UnauthenticatedRequestFactory{}.PutJson(config.GetActiveTarget(), path, query, body)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ag AuthenticatedRequester) PutJson(client *http.Client, config Config, path string, query string, body interface{}) ([]byte, error) {
	req, err := AuthenticatedRequestFactory{}.PutJson(config.GetActiveTarget(), path, query, body)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ug UnauthenticatedRequester) PatchJson(client *http.Client, config Config, path string, query string, body interface{}) ([]byte, error) {
	req, err := UnauthenticatedRequestFactory{}.PatchJson(config.GetActiveTarget(), path, query, body)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	return doAndRead(req, client, config)
}

func (ag AuthenticatedRequester) PatchJson(client *http.Client, config Config, path string, query string, body interface{}, extraHeaders map[string]string) ([]byte, error) {
	req, err := AuthenticatedRequestFactory{}.PatchJson(config.GetActiveTarget(), path, query, body)
	if err != nil {
		return []byte{}, err
	}
	addZoneSwitchHeader(req, &config)
	for k, v := range extraHeaders {
		req.Header.Add(k, v)
	}
	return doAndRead(req, client, config)
}
