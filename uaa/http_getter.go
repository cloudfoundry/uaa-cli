package uaa

import (
	"io/ioutil"
	"net/http"
	"github.com/jhamon/uaa-cli/utils"
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