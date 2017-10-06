package uaa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ScimGroupMember struct {
	Origin string `json:"origin,omitempty"`
	Type   string `json:"type,omitempty"`
	Value  string `json:"value,omitempty"`
}

type ScimGroup struct {
	ID          string            `json:"id,omitempty"`
	Meta        *ScimMetaInfo     `json:"meta,omitempty"`
	DisplayName string            `json:"displayName,omitempty"`
	ZoneID      string            `json:"zoneId,omitempty"`
	Description string            `json:"description,omitempty"`
	Members     []ScimGroupMember `json:"members,omitempty"`
	Schemas     []string          `json:"schemas,omitempty"`
}

type PaginatedGroupList struct {
	Resources    []ScimGroup `json:"resources"`
	StartIndex   int32       `json:"startIndex"`
	ItemsPerPage int32       `json:"itemsPerPage"`
	TotalResults int32       `json:"totalResults"`
	Schemas      []string    `json:"schemas"`
}

type GroupManager struct {
	HttpClient *http.Client
	Config     Config
}

type GroupMembership struct {
	Origin string `json:"origin"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}

func (gm GroupManager) AddMember(groupID, userID string) error {
	url := fmt.Sprintf("/Groups/%s/members", groupID)
	membership := GroupMembership{Origin: "uaa", Type: "USER", Value: userID}
	_, err := AuthenticatedRequester{}.PostJson(gm.HttpClient, gm.Config, url, "", membership)
	if err != nil {
		return err
	}

	return nil
}

func (gm GroupManager) Get(groupID string) (ScimGroup, error) {
	url := "/Groups/" + groupID
	bytes, err := AuthenticatedRequester{}.Get(gm.HttpClient, gm.Config, url, "")
	if err != nil {
		return ScimGroup{}, err
	}

	group := ScimGroup{}
	err = json.Unmarshal(bytes, &group)
	if err != nil {
		return ScimGroup{}, parseError(url, bytes)
	}

	return group, err
}

func (gm GroupManager) GetByName(name, attributes string) (ScimGroup, error) {
	if name == "" {
		return ScimGroup{}, errors.New("Group name may not be blank.")
	}

	filter := fmt.Sprintf(`displayName eq "%v"`, name)
	groups, err := gm.List(filter, "", attributes, "", 0, 0)
	if err != nil {
		return ScimGroup{}, err
	}
	if len(groups.Resources) == 0 {
		return ScimGroup{}, fmt.Errorf("Group %v not found.", name)
	}
	return groups.Resources[0], nil
}

func (gm GroupManager) List(filter, sortBy, attributes string, sortOrder ScimSortOrder, startIdx, count int) (PaginatedGroupList, error) {
	endpoint := "/Groups"

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

	bytes, err := AuthenticatedRequester{}.Get(gm.HttpClient, gm.Config, endpoint, query.Encode())
	if err != nil {
		return PaginatedGroupList{}, err
	}

	groupsResp := PaginatedGroupList{}
	err = json.Unmarshal(bytes, &groupsResp)
	if err != nil {
		return PaginatedGroupList{}, parseError(endpoint, bytes)
	}

	return groupsResp, err
}

func (gm GroupManager) Create(toCreate ScimGroup) (ScimGroup, error) {
	url := "/Groups"
	bytes, err := AuthenticatedRequester{}.PostJson(gm.HttpClient, gm.Config, url, "", toCreate)
	if err != nil {
		return ScimGroup{}, err
	}

	created := ScimGroup{}
	err = json.Unmarshal(bytes, &created)
	if err != nil {
		return ScimGroup{}, parseError(url, bytes)
	}

	return created, err
}

func (gm GroupManager) Update(toUpdate ScimGroup) (ScimGroup, error) {
	url := "/Groups"
	bytes, err := AuthenticatedRequester{}.PutJson(gm.HttpClient, gm.Config, url, "", toUpdate)
	if err != nil {
		return ScimGroup{}, err
	}

	updated := ScimGroup{}
	err = json.Unmarshal(bytes, &updated)
	if err != nil {
		return ScimGroup{}, parseError(url, bytes)
	}

	return updated, err
}

func (gm GroupManager) Delete(groupID string) (ScimGroup, error) {
	url := "/Groups/" + groupID
	bytes, err := AuthenticatedRequester{}.Delete(gm.HttpClient, gm.Config, url, "")
	if err != nil {
		return ScimGroup{}, err
	}

	deleted := ScimGroup{}
	err = json.Unmarshal(bytes, &deleted)
	if err != nil {
		return ScimGroup{}, parseError(url, bytes)
	}

	return deleted, err
}
