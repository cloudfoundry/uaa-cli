package uaa

import (
	"encoding/json"
	"net/http"
)

type ClientManager struct {
	HttpClient *http.Client
	Config     Config
}

type UaaClient struct {
	Scope                []string `json:"scope,omitempty"`
	ClientId             string   `json:"client_id,omitempty"`
	ClientSecret         string   `json:"client_secret,omitempty"`
	ResourceIds          []string `json:"resource_ids,omitempty"`
	AuthorizedGrantTypes []string `json:"authorized_grant_types,omitempty"`
	RedirectUri          []string `json:"redirect_uri,omitempty"`
	Autoapprove          []string `json:"autoapprove,omitempty"`
	Authorities          []string `json:"authorities,omitempty"`
	TokenSalt            string   `json:"token_salt,omitempty"`
	AllowedProviders     []string `json:"allowedproviders,omitempty"`
	DisplayName          string   `json:"name,omitempty"`
	LastModified         int64    `json:"lastModified,omitempty"`
	RequiredUserGroups   []string `json:"required_user_groups,omitempty"`
}

type PaginatedClientList struct {
	Resources []UaaClient
}

func (cm *ClientManager) Get(clientId string) (UaaClient, error) {
	url := "/oauth/clients/" + clientId
	bytes, err := AuthenticatedRequester{}.Get(cm.HttpClient, cm.Config, url, "")
	if err != nil {
		return UaaClient{}, err
	}

	uaaClient := UaaClient{}
	err = json.Unmarshal(bytes, &uaaClient)
	if err != nil {
		return UaaClient{}, parseError("", bytes)
	}

	return uaaClient, err
}

func (cm *ClientManager) Delete(clientId string) (UaaClient, error) {
	url := "/oauth/clients/" + clientId
	bytes, err := AuthenticatedRequester{}.Delete(cm.HttpClient, cm.Config, url, "")
	if err != nil {
		return UaaClient{}, err
	}

	uaaClient := UaaClient{}
	err = json.Unmarshal(bytes, &uaaClient)
	if err != nil {
		return UaaClient{}, parseError("", bytes)
	}

	return uaaClient, err
}

func (cm *ClientManager) Create(toCreate UaaClient) (UaaClient, error) {
	url := "/oauth/clients"
	bytes, err := AuthenticatedRequester{}.PostJson(cm.HttpClient, cm.Config, url, "", toCreate)
	if err != nil {
		return UaaClient{}, err
	}

	uaaClient := UaaClient{}
	err = json.Unmarshal(bytes, &uaaClient)
	if err != nil {
		return UaaClient{}, parseError("", bytes)
	}

	return uaaClient, err
}

func (cm *ClientManager) Update(toUpdate UaaClient) (UaaClient, error) {
	url := "/oauth/clients/" + toUpdate.ClientId
	bytes, err := AuthenticatedRequester{}.PutJson(cm.HttpClient, cm.Config, url, "", toUpdate)
	if err != nil {
		return UaaClient{}, err
	}

	uaaClient := UaaClient{}
	err = json.Unmarshal(bytes, &uaaClient)
	if err != nil {
		return UaaClient{}, parseError("", bytes)
	}

	return uaaClient, err
}

func (cm *ClientManager) List() ([]UaaClient, error) {
	// TODO: implement this
	return []UaaClient{}, nil
}
