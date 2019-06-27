package utils

import (
	"net/url"
	"path"
)

func BuildUrl(baseUrl, newPath string) (*url.URL, error) {
	newUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	newUrl.Path = path.Join(newUrl.Path, newPath)
	return newUrl, nil
}
