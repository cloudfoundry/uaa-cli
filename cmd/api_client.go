package cmd

import (
	"github.com/cloudfoundry-community/go-uaa"
	"net/http"
	"code.cloudfoundry.org/uaa-cli/config"
)

func GetAPIFromSavedTokenInContext() *uaa.API {
	config := GetSavedConfig()
	api, err := uaa.NewWithToken(config.GetActiveTarget().BaseUrl, config.ZoneSubdomain, config.GetActiveContext().Token, config.GetActiveTarget().SkipSSLValidation)

	if err != nil {
		panic(err)
	}

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
