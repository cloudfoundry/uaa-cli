package uaa

import "errors"

func requestError(url string) error {
	return errors.New("An unknown error occurred while calling " + url)
}

func parseError(url string, body []byte) error {
	errorMsg := "An unknown error occurred while parsing response from " + url + ". Response was " + string(body)
	return errors.New(errorMsg)
}

func unknownError() error {
	return errors.New("An unknown error occurred")
}
