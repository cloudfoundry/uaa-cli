package utils

import (
	"net/url"
)

func BuildUrl(baseUrl, path string) (*url.URL, error) {
	newUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	newUrl.Path = path
	return newUrl, nil
}
