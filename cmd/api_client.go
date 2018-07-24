package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/cloudfoundry-community/go-uaa"
	"net/http"
)

func GetAPIFromSavedTokenInContext() *uaa.API {
	config := GetSavedConfig()
	api, err := uaa.NewWithToken(config.GetActiveTarget().BaseUrl, config.ZoneSubdomain, config.GetActiveContext().Token)
	if err != nil {
		panic(err)
	}
	api.SkipSSLValidation = config.GetActiveTarget().SkipSSLValidation

	return api
}

func GetUnauthenticatedAPI() *uaa.API {
	config := GetSavedConfig()
	return GetUnauthenticatedAPIFromConfig(config)
}

func GetUnauthenticatedAPIFromConfig(cfg config.Config) *uaa.API {
	u, err := uaa.BuildTargetURL(cfg.GetActiveTarget().BaseUrl)
	if err != nil {
		panic(err)
	}

	client := &http.Client{Transport: http.DefaultTransport}
	return &uaa.API{
		UnauthenticatedClient: client,
		SkipSSLValidation:     cfg.GetActiveTarget().SkipSSLValidation,
		AuthenticatedClient:   nil,
		TargetURL:             u,
		ZoneID:                cfg.ZoneSubdomain,
	}
}
