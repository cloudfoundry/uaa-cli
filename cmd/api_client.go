package cmd

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/cloudfoundry-community/go-uaa"
)

func GetAPIFromSavedTokenInContext() *uaa.API {
	cfg := GetSavedConfig()
	token := cfg.GetActiveContext().Token
	api, err := uaa.New(
		cfg.GetActiveTarget().BaseUrl,
		uaa.WithToken(&token),
		uaa.WithZoneID(cfg.ZoneSubdomain),
		uaa.WithSkipSSLValidation(cfg.GetActiveTarget().SkipSSLValidation),
		uaa.WithVerbosity(verbose),
	)
	if err != nil {
		panic(err)
	}
	return api
}

func GetUnauthenticatedAPI() *uaa.API {
	cfg := GetSavedConfig()
	return GetUnauthenticatedAPIFromConfig(cfg)
}

func GetUnauthenticatedAPIFromConfig(cfg config.Config) *uaa.API {
	client := &http.Client{Transport: http.DefaultTransport}
	api, err := uaa.New(
		cfg.GetActiveTarget().BaseUrl,
		uaa.WithNoAuthentication(),
		uaa.WithZoneID(cfg.ZoneSubdomain),
		uaa.WithSkipSSLValidation(cfg.GetActiveTarget().SkipSSLValidation),
		uaa.WithVerbosity(verbose),
		uaa.WithClient(client),
	)
	if err != nil {
		panic(err)
	}
	return api
}
