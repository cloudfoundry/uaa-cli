package uaa

import (
	"encoding/json"
)

type Keys struct {
	Keys []JWK
}

func TokenKeys(context UaaContext) ([]JWK, error) {
	body, err := AuthenticatedGetter{}.Get(context, "/token_keys", "")
	if err != nil {
		key, err := TokenKey(context)
		return []JWK{key}, err
	}

	keys := Keys{}
	err = json.Unmarshal(body,&keys)
	if err != nil {
		return []JWK{}, parseError("/token_keys", body)
	}

	return keys.Keys, nil
}