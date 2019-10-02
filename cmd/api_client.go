package cmd

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/cloudfoundry-community/go-uaa"
)

func GetAPIFromSavedTokenInContext() *uaa.API {
	cfg := GetSavedConfig()
	api, err := uaa.NewWithToken(cfg.GetActiveTarget().BaseUrl, cfg.ZoneSubdomain, cfg.GetActiveContext().Token)
	if err != nil {
		panic(err)
	}
	api.WithSkipSSLValidation(cfg.GetActiveTarget().SkipSSLValidation)
	api.Verbose = verbose

	return api
}

func GetUnauthenticatedAPI() *uaa.API {
	cfg := GetSavedConfig()
	return GetUnauthenticatedAPIFromConfig(cfg)
}

func GetUnauthenticatedAPIFromConfig(cfg config.Config) *uaa.API {
	u, err := uaa.BuildTargetURL(cfg.GetActiveTarget().BaseUrl)
	if err != nil {
		panic(err)
	}

	client := &http.Client{Transport: http.DefaultTransport}

	api := &uaa.API{
		UnauthenticatedClient: client,
		AuthenticatedClient:   nil,
		TargetURL:             u,
		ZoneID:                cfg.ZoneSubdomain,
		Verbose:               verbose,
	}
	return api.WithSkipSSLValidation(cfg.GetActiveTarget().SkipSSLValidation)
}
