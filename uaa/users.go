package uaa

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type ScimMetaInfo struct {
	Version      int    `json:"version,omitempty"`
	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
}

type ScimUserName struct {
	FamilyName string `json:"familyName,omitempty"`
	GivenName  string `json:"givenName,omitempty"`
}

type ScimUserEmail struct {
	Value   string `json:"value,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

type ScimGroup struct {
	Value   string `json:"value,omitempty"`
	Display string `json:"display,omitempty"`
	Type    string `json:"type,omitempty"`
}

type Approval struct {
	UserId        string `json:"userId,omitempty"`
	ClientId      string `json:"clientId,omitempty"`
	Scope         string `json:"scope,omitempty"`
	Status        string `json:"status,omitempty"`
	LastUpdatedAt string `json:"lastUpdatedAt,omitempty"`
	ExpiresAt     string `json:"expiresAt,omitempty"`
}

type PhoneNumber struct {
	Value string `json:"value"`
}

type ScimUser struct {
	Id                   string          `json:"id,omitempty"`
	ExternalId           string          `json:"externalId,omitempty"`
	Meta                 *ScimMetaInfo   `json:"meta,omitempty"`
	Username             string          `json:"userName,omitempty"`
	Name                 *ScimUserName   `json:"name,omitempty"`
	Emails               []ScimUserEmail `json:"emails,omitempty"`
	Groups               []ScimGroup     `json:"groups,omitempty"`
	Approvals            []Approval      `json:"approvals,omitempty"`
	PhoneNumbers         []PhoneNumber   `json:"phoneNumbers,omitempty"`
	Active               bool            `json:"active,omitempty"`
	Verified             bool            `json:"verified,omitempty"`
	Origin               string          `json:"origin,omitempty"`
	ZoneId               string          `json:"zoneId,omitempty"`
	PasswordLastModified string          `json:"passwordLastModified,omitempty"`
	PreviousLogonTime    int             `json:"previousLogonTime,omitempty"`
	LastLogonTime        int             `json:"lastLogonTime,omitempty"`
	Schemas              []string        `json:"schemas,omitempty"`
}

type Crud interface {
	Get(string) (ScimUser, error)
	List(string, string, string, ScimSortOrder, int, int) (PaginatedUserList, error)
}

type UserManager struct {
	HttpClient *http.Client
	Config     Config
}

type PaginatedUserList struct {
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

func (um UserManager) List(filter, sortBy, attributes string, sortOrder ScimSortOrder, startIdx, count int) (PaginatedUserList, error) {
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
		return PaginatedUserList{}, err
	}

	usersResp := PaginatedUserList{}
	err = json.Unmarshal(bytes, &usersResp)
	if err != nil {
		return PaginatedUserList{}, parseError(endpoint, bytes)
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
func (tc TestUserCrud) List(filter, sortBy, attributes string, sortOrder ScimSortOrder, startIdx, count int) (PaginatedUserList, error) {
	return PaginatedUserList{}, nil
}

func NewTestUserCrud() TestUserCrud {
	tc := TestUserCrud{}
	tc.CallData = make(map[string]string)
	return tc
}
