package uaa

import (
	"net/http"
	"encoding/json"
)

type ClientManager struct {
	HttpClient *http.Client
	Config Config
}

type UaaClient struct {
	Scope []string `json:"scope"`
	ClientId string `json:"client_id"`
	ResourceIds []string `json:"resource_ids"`
	AuthorizedGrantTypes []string `json:"authorized_grant_types"`
	RedirectUri []string `json:"redirect_uri"`
	Autoapprove []string `json:"autoapprove"`
	Authorities []string `json:"authorities"`
	TokenSalt string `json:"token_salt"`
	AllowedProviders []string `json:"allowedproviders"`
	Name string `json:"name"`
	LastModified int64 `json:"lastModified"`
	RequiredUserGroups []string `json:"required_user_groups"`
}

func (cm *ClientManager) Get(clientId string) (UaaClient, error) {
	url := "/oauth/clients/" + clientId
	bytes, err := AuthenticatedRequester{}.GetBytes(cm.HttpClient, cm.Config, url, "")
	if err != nil {
		return UaaClient{}, err
	}

	uaaClient := UaaClient{}
	err = json.Unmarshal(bytes,&uaaClient)
	if err != nil {
		return UaaClient{}, parseError("", bytes)
	}

	return uaaClient, err
}