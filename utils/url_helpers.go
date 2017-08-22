package utils

import "net/url"

func BuildUrl(baseUrl, path string) string {
	newUrl, err := url.Parse(baseUrl)
	if err != nil {
		return ""
	}

	newUrl.Path = path
	return newUrl.String()
}
