package uaa

import (
	"code.cloudfoundry.org/uaa-cli/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type ClientManager struct {
	HttpClient *http.Client
	Config     Config
}

type UaaClient struct {
	ClientId             string   `json:"client_id,omitempty"`
	ClientSecret         string   `json:"client_secret,omitempty"`
	Scope                []string `json:"scope,omitempty"`
	ResourceIds          []string `json:"resource_ids,omitempty"`
	AuthorizedGrantTypes []string `json:"authorized_grant_types,omitempty"`
	RedirectUri          []string `json:"redirect_uri,omitempty"`
	Authorities          []string `json:"authorities,omitempty"`
	TokenSalt            string   `json:"token_salt,omitempty"`
	AllowedProviders     []string `json:"allowedproviders,omitempty"`
	DisplayName          string   `json:"name,omitempty"`
	LastModified         int64    `json:"lastModified,omitempty"`
	RequiredUserGroups   []string `json:"required_user_groups,omitempty"`
}

func errorMissingValueForGrantType(value string, grantType GrantType) error {
	return errors.New(fmt.Sprintf("%v must be specified for %v grant type.", value, grantType))
}

func errorMissingValue(value string) error {
	return errors.New(fmt.Sprintf("%v must be specified in the client definition.", value))
}

func requireRedirectUriForGrantType(c *UaaClient, grantType GrantType) error {
	if utils.Contains(c.AuthorizedGrantTypes, string(grantType)) {
		if len(c.RedirectUri) == 0 {
			return errorMissingValueForGrantType("redirect_uri", grantType)
		}
	}
	return nil
}

func requireClientSecretForGrantType(c *UaaClient, grantType GrantType) error {
	if utils.Contains(c.AuthorizedGrantTypes, string(grantType)) {
		if c.ClientSecret == "" {
			return errorMissingValueForGrantType("client_secret", grantType)
		}
	}
	return nil
}

func knownGrantTypesStr() string {
	grantTypeStrings := []string{}
	KNOWN_GRANT_TYPES := []GrantType{AUTHCODE, IMPLICIT, PASSWORD, CLIENT_CREDENTIALS}
	for _, grant := range KNOWN_GRANT_TYPES {
		grantTypeStrings = append(grantTypeStrings, string(grant))
	}

	return "[" + strings.Join(grantTypeStrings, ", ") + "]"
}

func (c *UaaClient) PreCreateValidation() error {
	if len(c.AuthorizedGrantTypes) == 0 {
		return errors.New(fmt.Sprintf("Grant type must be one of %v", knownGrantTypesStr()))
	}

	if c.ClientId == "" {
		return errorMissingValue("client_id")
	}

	if err := requireRedirectUriForGrantType(c, AUTHCODE); err != nil {
		return err
	}
	if err := requireClientSecretForGrantType(c, AUTHCODE); err != nil {
		return err
	}

	if err := requireClientSecretForGrantType(c, PASSWORD); err != nil {
		return err
	}

	if err := requireClientSecretForGrantType(c, CLIENT_CREDENTIALS); err != nil {
		return err
	}

	if err := requireRedirectUriForGrantType(c, IMPLICIT); err != nil {
		return err
	}

	return nil
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
