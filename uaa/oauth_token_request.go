package uaa

import (
	"net/http"
	"encoding/json"
)

type GrantType string

type ClientCredentialsClient struct {
	ClientId string
	ClientSecret string
}

func (cc ClientCredentialsClient) RequestToken(httpClient *http.Client, config Config, format TokenFormat) (TokenResponse, error) {
	body := map[string]string{
		"client_id": cc.ClientId,
		"client_secret": cc.ClientSecret,
		"grant_type": "client_credentials",
		"token_format": string(format),
		"response_type": "token",
	}

	bytes, err := UnauthenticatedRequester{}.PostBytes(httpClient, config, "/oauth/token", "", body)
	if err != nil {
		return TokenResponse{}, err
	}

	tokenResponse := TokenResponse{}
	err = json.Unmarshal(bytes,&tokenResponse)
	if err != nil {
		return TokenResponse{}, parseError("/oauth/token", bytes)
	}

	return tokenResponse, nil
}

type TokenFormat string
const (
	OPAQUE = TokenFormat("opaque")
	JWT = TokenFormat("jwt")
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int32 `json:"expires_in"`
	Scope string `json:"scope"`
	JTI string `json:"jti"`
}