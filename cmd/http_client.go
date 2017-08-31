package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"crypto/tls"
	"net/http"
	"time"
)

func GetHttpClient() *http.Client {
	return GetHttpClientWithConfig(GetSavedConfig())
}

// This should really only be called directly by the target
// command, as it wants to build an http client before saving
// the new target.
func GetHttpClientWithConfig(config uaa.Config) *http.Client {
	var client = &http.Client{
		Timeout: 5 * time.Second,
	}

	if config.GetActiveTarget().SkipSSLValidation {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	return client
}

func GetSavedConfig() uaa.Config {
	c := config.ReadConfig()
	c.Trace = trace
	c.ZoneSubdomain = zoneSubdomain
	return c
}
