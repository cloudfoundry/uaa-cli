package uaa

import (
	"net/http"
	"encoding/json"
)

type ClientCredentialsClient struct {
	ClientId string
	ClientSecret string
}

func postToOAuthToken(httpClient *http.Client, config Config, body map[string]string) (TokenResponse, error) {
	bytes, err := UnauthenticatedRequester{}.PostForm(httpClient, config, "/oauth/token", "", body)
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

func (cc ClientCredentialsClient) RequestToken(httpClient *http.Client, config Config, format TokenFormat) (TokenResponse, error) {
	body := map[string]string{
		"grant_type": "client_credentials",
		"client_id": cc.ClientId,
		"client_secret": cc.ClientSecret,
		"token_format": string(format),
		"response_type": "token",
	}

	return postToOAuthToken(httpClient, config, body)
}

type ResourceOwnerPasswordClient struct {
	ClientId string
	ClientSecret string
	Username string
	Password string
}

func (rop ResourceOwnerPasswordClient) RequestToken(httpClient *http.Client, config Config, format TokenFormat) (TokenResponse, error) {
	body := map[string]string{
		"grant_type": "password",
		"client_id": rop.ClientId,
		"client_secret": rop.ClientSecret,
		"username": rop.Username,
		"password": rop.Password,
		"token_format": string(format),
		"response_type": "token",
	}

	return postToOAuthToken(httpClient, config, body)
}

type TokenFormat string
const (
	OPAQUE = TokenFormat("opaque")
	JWT = TokenFormat("jwt")
)

type GrantType string
const (
	CLIENT_CREDENTIALS = GrantType("client_credentials")
	PASSWORD = GrantType("password")
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int32 `json:"expires_in"`
	Scope string `json:"scope"`
	JTI string `json:"jti"`
}