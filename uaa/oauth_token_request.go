package uaa

import (
	"encoding/json"
	"net/http"
)

func postToOAuthToken(httpClient *http.Client, config Config, body map[string]string) (TokenResponse, error) {
	bytes, err := UnauthenticatedRequester{}.PostForm(httpClient, config, "/oauth/token", "", body)
	if err != nil {
		return TokenResponse{}, err
	}

	tokenResponse := TokenResponse{}
	err = json.Unmarshal(bytes, &tokenResponse)
	if err != nil {
		return TokenResponse{}, parseError("/oauth/token", bytes)
	}

	return tokenResponse, nil
}

type ClientCredentialsClient struct {
	ClientId     string
	ClientSecret string
}

func (cc ClientCredentialsClient) RequestToken(httpClient *http.Client, config Config, format TokenFormat) (TokenResponse, error) {
	body := map[string]string{
		"grant_type":    string(CLIENT_CREDENTIALS),
		"client_id":     cc.ClientId,
		"client_secret": cc.ClientSecret,
		"token_format":  string(format),
		"response_type": "token",
	}

	return postToOAuthToken(httpClient, config, body)
}

type ResourceOwnerPasswordClient struct {
	ClientId     string
	ClientSecret string
	Username     string
	Password     string
}

func (rop ResourceOwnerPasswordClient) RequestToken(httpClient *http.Client, config Config, format TokenFormat) (TokenResponse, error) {
	body := map[string]string{
		"grant_type":    string(PASSWORD),
		"client_id":     rop.ClientId,
		"client_secret": rop.ClientSecret,
		"username":      rop.Username,
		"password":      rop.Password,
		"token_format":  string(format),
		"response_type": "token",
	}

	return postToOAuthToken(httpClient, config, body)
}

type AuthorizationCodeClient struct {
	ClientId     string
	ClientSecret string
}

func (acc AuthorizationCodeClient) RequestToken(httpClient *http.Client, config Config, format TokenFormat, code string, redirectUri string) (TokenResponse, error) {
	body := map[string]string{
		"grant_type":    string(AUTHCODE),
		"client_id":     acc.ClientId,
		"client_secret": acc.ClientSecret,
		"token_format":  string(format),
		"response_type": "token",
		"redirect_uri":  redirectUri,
		"code":          code,
	}

	return postToOAuthToken(httpClient, config, body)
}

type RefreshTokenClient struct {
	ClientId     string
	ClientSecret string
}

func (rc RefreshTokenClient) RequestToken(httpClient *http.Client, config Config, format TokenFormat, refreshToken string) (TokenResponse, error) {
	body := map[string]string{
		"grant_type":    string(REFRESH_TOKEN),
		"refresh_token": refreshToken,
		"client_id":     rc.ClientId,
		"client_secret": rc.ClientSecret,
		"token_format":  string(format),
		"response_type": "token",
	}

	return postToOAuthToken(httpClient, config, body)
}

type TokenFormat string

const (
	OPAQUE = TokenFormat("opaque")
	JWT    = TokenFormat("jwt")
)

type GrantType string

const (
	REFRESH_TOKEN      = GrantType("refresh_token")
	AUTHCODE           = GrantType("authorization_code")
	IMPLICIT           = GrantType("implicit")
	PASSWORD           = GrantType("password")
	CLIENT_CREDENTIALS = GrantType("client_credentials")
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int32  `json:"expires_in"`
	Scope        string `json:"scope"`
	JTI          string `json:"jti"`
}
