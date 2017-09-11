package uaa

import (
	"net/http"
	"encoding/json"
)

type ScimMetaInfo struct {
	Version      int    `json:"version"`
	Created      string `json:"created"`
	LastModified string `json:"lastModified"`
}

type ScimUserName struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

type ScimUserEmail struct {
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}

type ScimGroup struct {
	Value   string `json:"value"`
	Display string `json:"display"`
	Type    string `json:"type"`
}

type Approval struct {
	UserId        string `json:"userId"`
	ClientId      string `json:"clientId"`
	Scope         string `json:"scope"`
	Status        string `json:"status"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
	ExpiresAt     string `json:"expiresAt"`
}

type PhoneNumber struct {
	Value string `json:"value"`
}

type ScimUser struct {
	Id                   string          `json:"id"`
	ExternalId           string          `json:"externalId"`
	Meta                 ScimMetaInfo    `json:"meta"`
	Username             string          `json:"userName"`
	Name                 ScimUserName    `json:"name"`
	Emails               []ScimUserEmail `json:"emails"`
	Groups               []ScimGroup     `json:"groups"`
	Approvals            []Approval      `json:"approvals"`
	PhoneNumbers         []PhoneNumber   `json:"phoneNumbers"`
	Active               bool            `json:"active"`
	Verified             bool            `json:"verified"`
	Origin               string          `json:"origin"`
	ZoneId               string          `json:"zoneId"`
	PasswordLastModified string          `json:"passwordLastModified"`
	PreviousLogonTime    int             `json:"previousLogonTime"`
	LastLogonTime        int             `json:"lastLogonTime"`
	Schemas              []string        `json:"schemas"`
}

type UserManager struct {
	HttpClient *http.Client
	Config     Config
}

func (um *UserManager) Get(userId string) (ScimUser, error) {
	url := "/Users/" + userId
	bytes, err := AuthenticatedRequester{}.Get(um.HttpClient, um.Config, url, "")
	if err != nil {
		return ScimUser{}, err
	}

	user := ScimUser{}
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return ScimUser{}, parseError(url, bytes)
	}

	return user, err
}


