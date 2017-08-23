package uaa

import (
	"net/http"
	"github.com/jhamon/guac/utils"
)

type UaaHealthStatus string

const (
	OK    = UaaHealthStatus("ok")
	ERROR = UaaHealthStatus("health_error")
)

func Health(context UaaContext) (UaaHealthStatus, error) {
	url, err := utils.BuildUrl(context.BaseUrl, "healthz")
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return "", nil
	}

	if resp.StatusCode == 200 {
		return OK, nil
	} else {
		return ERROR, nil
	}
}
