package uaa

import (
	"net/http"
	"github.com/jhamon/uaa-cli/utils"
)

type UaaHealthStatus string

const (
	OK    = UaaHealthStatus("ok")
	ERROR = UaaHealthStatus("health_error")
)

func Health(target Target) (UaaHealthStatus, error) {
	url, err := utils.BuildUrl(target.BaseUrl, "healthz")
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
