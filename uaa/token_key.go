package uaa

import (
	"encoding/json"
	"net/http"
)

type JWK struct {
	Kty   string `json:"kty"`
	E     string `json:"e,omitempty"`
	Use   string `json:"use"`
	Kid   string `json:"kid"`
	Alg   string `json:"alg"`
	Value string `json:"value"`
	N     string `json:"n,omitempty"`
}

func TokenKey(client *http.Client, config Config) (JWK, error) {
	body, err := UnauthenticatedRequester{}.Get(client, config, "token_key", "")
	if err != nil {
		return JWK{}, err
	}

	key := JWK{}
	err = json.Unmarshal(body, &key)
	if err != nil {
		return JWK{}, parseError("/token_key", body)
	}

	return key, nil
}
