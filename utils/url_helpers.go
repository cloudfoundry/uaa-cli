package utils

import "net/url"

func BuildUrl(baseUrl, path string) *url.URL {
	newUrl, _ := url.Parse(baseUrl)

	newUrl.Path = path
	return newUrl
}
