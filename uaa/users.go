package uaa

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
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

type Crud interface {
	Get(string) (ScimUser, error)
	List(string) ([]ScimUser, error)
}

type UserManager struct {
	HttpClient *http.Client
	Config     Config
}

type PaginatedUserResponse struct {
	Resources    []ScimUser `json:"resources"`
	StartIndex   int32      `json:"startIndex"`
	ItemsPerPage int32      `json:"itemsPerPage"`
	TotalResults int32      `json:"totalResults"`
	Schemas      []string   `json:"schemas"`
}

func (um UserManager) Get(userId string) (ScimUser, error) {
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

type ScimSortOrder string

const (
	SORT_ASCENDING  = ScimSortOrder("ascending")
	SORT_DESCENDING = ScimSortOrder("descending")
)

func (um UserManager) List(filter, sortBy, attributes string, sortOrder ScimSortOrder, startIdx, count int) (PaginatedUserResponse, error) {
	endpoint := "/Users"

	query := url.Values{}
	if filter != "" {
		query.Add("filter", filter)
	}
	if attributes != "" {
		query.Add("attributes", attributes)
	}
	if sortBy != "" {
		query.Add("sortBy", sortBy)
	}
	if count != 0 {
		query.Add("count", strconv.Itoa(count))
	}
	if startIdx != 0 {
		query.Add("startIndex", strconv.Itoa(startIdx))
	}
	if sortOrder != "" {
		query.Add("sortOrder", string(sortOrder))
	}

	bytes, err := AuthenticatedRequester{}.Get(um.HttpClient, um.Config, endpoint, query.Encode())
	if err != nil {
		return PaginatedUserResponse{}, err
	}

	usersResp := PaginatedUserResponse{}
	err = json.Unmarshal(bytes, &usersResp)
	if err != nil {
		return PaginatedUserResponse{}, parseError(endpoint, bytes)
	}

	return usersResp, err
}

type TestUserCrud struct {
	CallData map[string]string
}

func (tc TestUserCrud) Get(id string) (ScimUser, error) {
	tc.CallData["GetId"] = id
	return ScimUser{Id: id}, nil
}
func NewTestUserCrud() TestUserCrud {
	tc := TestUserCrud{}
	tc.CallData = make(map[string]string)
	return tc
}
