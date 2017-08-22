package info

import (
	"net/http"
	"github.com/jhamon/uaalib/utils"
)

type UaaClient struct {
	BaseUrl string
}

type UaaHealthStatus string

const (
	OK    = UaaHealthStatus("ok")
	ERROR = UaaHealthStatus("health_error")
)

func Health(client UaaClient) UaaHealthStatus {
		resp, _ := http.Get(utils.BuildUrl(client.BaseUrl, "healthz"))

	if resp.StatusCode == 200 {
		return OK
	} else {
		return ERROR
	}
}
