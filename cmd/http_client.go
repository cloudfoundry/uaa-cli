package cmd

import (
	"net/http"
	"crypto/tls"
	"time"
)

func GetHttpClient() *http.Client {
	var client = &http.Client{
		Timeout: 5 * time.Second,
	}

	if (skipSSLValidation) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	return client
}
