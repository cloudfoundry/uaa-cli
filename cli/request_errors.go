package cli

import (
	"github.com/cloudfoundry-community/go-uaa"
	"golang.org/x/oauth2"
)

func RequestErrorFromOauthError(err error) error {
	oauthErrorResponse, isRetrieveError := err.(*oauth2.RetrieveError)
	if isRetrieveError {
		tokenUrl := oauthErrorResponse.Response.Request.URL.String()
		return uaa.RequestError{Url: tokenUrl, ErrorResponse: oauthErrorResponse.Body}
	}
	return err
}
