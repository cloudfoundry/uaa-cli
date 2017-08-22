package info

import (
	"net/http"
	"net/url"
)

type UaaClient struct {
	BaseUrl string
}

type UaaHealthStatus string

const (
	OK = UaaHealthStatus("ok")
	ERROR = UaaHealthStatus("health_error")
)

func buildUrl(baseUrl, path string) string {
	newUrl, err := url.Parse(baseUrl)
	if (err != nil) {
		return ""
	}

	newUrl.Path = path
	return newUrl.String()
}

func Health(client UaaClient) UaaHealthStatus {
	resp, _ := http.Get(buildUrl(client.BaseUrl, "healthz"))

	if (resp.StatusCode == 200) {
		return OK
	} else {
		return ERROR
	}
}
