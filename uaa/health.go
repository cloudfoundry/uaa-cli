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

func Health(context UaaContext) UaaHealthStatus {
		resp, _ := http.Get(utils.BuildUrl(context.BaseUrl, "healthz").String())

	if resp.StatusCode == 200 {
		return OK
	} else {
		return ERROR
	}
}
