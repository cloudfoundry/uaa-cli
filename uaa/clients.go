package uaa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ClientManager struct {
	HttpClient *http.Client
	Config     Config
}

type UaaClient struct {
	Scope                []string
	ClientId             string
	ClientSecret         string
	ResourceIds          []string
	AuthorizedGrantTypes []string
	RedirectUri          []string
	Authorities          []string
	AutoapprovedScopes   []string
	AutoapproveAll       bool
	TokenSalt            string
	AllowedProviders     []string
	DisplayName          string
	LastModified         int64
	RequiredUserGroups   []string
}

// This json-parsing insanity is necessary because the UAA will accept
// either a boolean or an array of strings in the autoapprove field
func (uc UaaClient) MarshalJSON() ([]byte, error) {
	if uc.AutoapproveAll {
		toJsonify := &struct {
			AutoapproveAll       bool     `json:"autoapprove,omitempty"`
			Scope                []string `json:"scope,omitempty"`
			ClientId             string   `json:"client_id,omitempty"`
			ClientSecret         string   `json:"client_secret,omitempty"`
			ResourceIds          []string `json:"resource_ids,omitempty"`
			AuthorizedGrantTypes []string `json:"authorized_grant_types,omitempty"`
			RedirectUri          []string `json:"redirect_uri,omitempty"`
			Authorities          []string `json:"authorities,omitempty"`
			TokenSalt            string   `json:"token_salt,omitempty"`
			AllowedProviders     []string `json:"allowedproviders,omitempty"`
			DisplayName          string   `json:"name,omitempty"`
			LastModified         int64    `json:"lastModified,omitempty"`
			RequiredUserGroups   []string `json:"required_user_groups,omitempty"`
		}{
			AutoapproveAll:       uc.AutoapproveAll,
			Scope:                uc.Scope,
			ClientId:             uc.ClientId,
			ClientSecret:         uc.ClientSecret,
			ResourceIds:          uc.ResourceIds,
			AuthorizedGrantTypes: uc.AuthorizedGrantTypes,
			RedirectUri:          uc.RedirectUri,
			Authorities:          uc.Authorities,
			TokenSalt:            uc.TokenSalt,
			AllowedProviders:     uc.AllowedProviders,
			DisplayName:          uc.DisplayName,
			LastModified:         uc.LastModified,
			RequiredUserGroups:   uc.RequiredUserGroups,
		}
		return json.Marshal(toJsonify)
	}

	toJsonify := &struct {
		AutoapprovedScopes   []string `json:"autoapprove,omitempty"`
		Scope                []string `json:"scope,omitempty"`
		ClientId             string   `json:"client_id,omitempty"`
		ClientSecret         string   `json:"client_secret,omitempty"`
		ResourceIds          []string `json:"resource_ids,omitempty"`
		AuthorizedGrantTypes []string `json:"authorized_grant_types,omitempty"`
		RedirectUri          []string `json:"redirect_uri,omitempty"`
		Authorities          []string `json:"authorities,omitempty"`
		TokenSalt            string   `json:"token_salt,omitempty"`
		AllowedProviders     []string `json:"allowedproviders,omitempty"`
		DisplayName          string   `json:"name,omitempty"`
		LastModified         int64    `json:"lastModified,omitempty"`
		RequiredUserGroups   []string `json:"required_user_groups,omitempty"`
	}{
		AutoapprovedScopes:   uc.AutoapprovedScopes,
		Scope:                uc.Scope,
		ClientId:             uc.ClientId,
		ClientSecret:         uc.ClientSecret,
		ResourceIds:          uc.ResourceIds,
		AuthorizedGrantTypes: uc.AuthorizedGrantTypes,
		RedirectUri:          uc.RedirectUri,
		Authorities:          uc.Authorities,
		TokenSalt:            uc.TokenSalt,
		AllowedProviders:     uc.AllowedProviders,
		DisplayName:          uc.DisplayName,
		LastModified:         uc.LastModified,
		RequiredUserGroups:   uc.RequiredUserGroups,
	}
	return json.Marshal(toJsonify)
}

// []byte to UaaClient
func (uc *UaaClient) UnmarshalJSON(data []byte) error {
	clientWithBoolAutoapprove := &struct {
		AutoapproveAll       bool     `json:"autoapprove,omitempty"`
		Scope                []string `json:"scope,omitempty"`
		ClientId             string   `json:"client_id,omitempty"`
		ClientSecret         string   `json:"client_secret,omitempty"`
		ResourceIds          []string `json:"resource_ids,omitempty"`
		AuthorizedGrantTypes []string `json:"authorized_grant_types,omitempty"`
		RedirectUri          []string `json:"redirect_uri,omitempty"`
		Authorities          []string `json:"authorities,omitempty"`
		TokenSalt            string   `json:"token_salt,omitempty"`
		AllowedProviders     []string `json:"allowedproviders,omitempty"`
		DisplayName          string   `json:"name,omitempty"`
		LastModified         int64    `json:"lastModified,omitempty"`
		RequiredUserGroups   []string `json:"required_user_groups,omitempty"`
	}{}

	err := json.Unmarshal(data, &clientWithBoolAutoapprove)
	if err == nil {
		uc.AutoapproveAll = clientWithBoolAutoapprove.AutoapproveAll
		uc.Scope = clientWithBoolAutoapprove.Scope
		uc.ClientId = clientWithBoolAutoapprove.ClientId
		uc.ClientSecret = clientWithBoolAutoapprove.ClientSecret
		uc.ResourceIds = clientWithBoolAutoapprove.ResourceIds
		uc.AuthorizedGrantTypes = clientWithBoolAutoapprove.AuthorizedGrantTypes
		uc.RedirectUri = clientWithBoolAutoapprove.RedirectUri
		uc.Authorities = clientWithBoolAutoapprove.Authorities
		uc.TokenSalt = clientWithBoolAutoapprove.TokenSalt
		uc.AllowedProviders = clientWithBoolAutoapprove.AllowedProviders
		uc.DisplayName = clientWithBoolAutoapprove.DisplayName
		uc.LastModified = clientWithBoolAutoapprove.LastModified
		uc.RequiredUserGroups = clientWithBoolAutoapprove.RequiredUserGroups
		return nil
	}

	clientWithStrArrAutoapprove := &struct {
		AutoapprovedScopes   []string `json:"autoapprove,omitempty"`
		Scope                []string `json:"scope,omitempty"`
		ClientId             string   `json:"client_id,omitempty"`
		ClientSecret         string   `json:"client_secret,omitempty"`
		ResourceIds          []string `json:"resource_ids,omitempty"`
		AuthorizedGrantTypes []string `json:"authorized_grant_types,omitempty"`
		RedirectUri          []string `json:"redirect_uri,omitempty"`
		Authorities          []string `json:"authorities,omitempty"`
		TokenSalt            string   `json:"token_salt,omitempty"`
		AllowedProviders     []string `json:"allowedproviders,omitempty"`
		DisplayName          string   `json:"name,omitempty"`
		LastModified         int64    `json:"lastModified,omitempty"`
		RequiredUserGroups   []string `json:"required_user_groups,omitempty"`
	}{}

	err = json.Unmarshal(data, &clientWithStrArrAutoapprove)
	if err == nil {
		uc.AutoapprovedScopes = clientWithStrArrAutoapprove.AutoapprovedScopes
		uc.Scope = clientWithStrArrAutoapprove.Scope
		uc.ClientId = clientWithStrArrAutoapprove.ClientId
		uc.ClientSecret = clientWithStrArrAutoapprove.ClientSecret
		uc.ResourceIds = clientWithStrArrAutoapprove.ResourceIds
		uc.AuthorizedGrantTypes = clientWithStrArrAutoapprove.AuthorizedGrantTypes
		uc.RedirectUri = clientWithStrArrAutoapprove.RedirectUri
		uc.Authorities = clientWithStrArrAutoapprove.Authorities
		uc.TokenSalt = clientWithStrArrAutoapprove.TokenSalt
		uc.AllowedProviders = clientWithStrArrAutoapprove.AllowedProviders
		uc.DisplayName = clientWithStrArrAutoapprove.DisplayName
		uc.LastModified = clientWithStrArrAutoapprove.LastModified
		uc.RequiredUserGroups = clientWithStrArrAutoapprove.RequiredUserGroups
		return nil
	}

	return errors.New("Could not Unmarshal client json: " + string(data) + " " + err.Error())
}

type changeSecretBody struct {
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"secret,omitempty"`
}

type PaginatedClientList struct {
	Resources    []UaaClient `json:"resources"`
	StartIndex   int         `json:"startIndex"`
	ItemsPerPage int         `json:"itemsPerPage"`
	TotalResults int         `json:"totalResults"`
	Schemas      []string    `json:"schemas"`
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

func (cm *ClientManager) ChangeSecret(clientId string, newSecret string) error {
	url := "/oauth/clients/" + clientId + "/secret"
	body := changeSecretBody{ClientId: clientId, ClientSecret: newSecret}
	_, err := AuthenticatedRequester{}.PutJson(cm.HttpClient, cm.Config, url, "", body)
	return err
}

func getResultPage(cm *ClientManager, startIndex, count int) (PaginatedClientList, error) {
	query := fmt.Sprintf("startIndex=%v&count=%v", startIndex, count)
	if startIndex == 0 {
		query = ""
	}

	bytes, err := AuthenticatedRequester{}.Get(cm.HttpClient, cm.Config, "/oauth/clients", query)
	if err != nil {
		return PaginatedClientList{}, err
	}

	clientList := PaginatedClientList{}
	err = json.Unmarshal(bytes, &clientList)
	if err != nil {
		return PaginatedClientList{}, parseError("/oauth/clients", bytes)
	}
	return clientList, nil
}

func (cm *ClientManager) List() ([]UaaClient, error) {
	results, err := getResultPage(cm, 0, 0)
	if err != nil {
		return []UaaClient{}, err
	}

	clientList := results.Resources
	startIndex, count := results.StartIndex, results.ItemsPerPage
	for results.TotalResults > len(clientList) {
		startIndex += count
		newResults, err := getResultPage(cm, startIndex, count)
		if err != nil {
			return []UaaClient{}, err
		}
		clientList = append(clientList, newResults.Resources...)
	}

	return clientList, nil
}
