package uaa

import (
	"encoding/json"
	"net/http"
)

type Keys struct {
	Keys []JWK
}

func TokenKeys(client *http.Client, context UaaContext) ([]JWK, error) {
	body, err := AuthenticatedGetter{}.GetBytes(client, context, "/token_keys", "")
	if err != nil {
		key, err := TokenKey(client, context)
		return []JWK{key}, err
	}

	keys := Keys{}
	err = json.Unmarshal(body,&keys)
	if err != nil {
		return []JWK{}, parseError("/token_keys", body)
	}

	return keys.Keys, nil
}