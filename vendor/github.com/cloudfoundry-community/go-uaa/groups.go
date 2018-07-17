package uaa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// GroupsEndpoint is the path to the groups resource.
const GroupsEndpoint string = "/Groups"

// paginatedGroupList is the response from the API for a single page of groups.
type paginatedGroupList struct {
	Page
	Resources []Group  `json:"resources"`
	Schemas   []string `json:"schemas"`
}

// GroupMember is a user or a group.
type GroupMember struct {
	Origin string `json:"origin,omitempty"`
	Type   string `json:"type,omitempty"`
	Value  string `json:"value,omitempty"`
}

// Group is a container for users and groups.
type Group struct {
	ID          string        `json:"id,omitempty"`
	Meta        *Meta         `json:"meta,omitempty"`
	DisplayName string        `json:"displayName,omitempty"`
	ZoneID      string        `json:"zoneId,omitempty"`
	Description string        `json:"description,omitempty"`
	Members     []GroupMember `json:"members,omitempty"`
	Schemas     []string      `json:"schemas,omitempty"`
}

// Identifier returns the field used to uniquely identify a Group.
func (g Group) Identifier() string {
	return g.ID
}

// AddGroupMember adds the entity with the given memberID to the group with the
// given ID. If no entityType is supplied, the entityType (which can be "USER"
// or "GROUP") will be "USER". If no origin is supplied, the origin will be
// "uaa".
func (a *API) AddGroupMember(groupID string, memberID string, entityType string, origin string) error {
	u := urlWithPath(*a.TargetURL, fmt.Sprintf("%s/%s/members", GroupsEndpoint, groupID))
	if origin == "" {
		origin = "uaa"
	}
	if entityType == "" {
		entityType = "USER"
	}
	membership := GroupMember{Origin: origin, Type: entityType, Value: memberID}
	j, err := json.Marshal(membership)
	if err != nil {
		return err
	}
	err = a.doJSON(http.MethodPost, &u, bytes.NewBuffer([]byte(j)), nil, true)
	if err != nil {
		return err
	}
	return nil
}

// RemoveGroupMember removes the entity with the given memberID from the group
// with the given ID. If no entityType is supplied, the entityType (which can be
// "USER" or "GROUP") will be "USER". If no origin is supplied, the origin will
// be "uaa".
func (a *API) RemoveGroupMember(groupID string, memberID string, entityType string, origin string) error {
	u := urlWithPath(*a.TargetURL, fmt.Sprintf("%s/%s/members/%s", GroupsEndpoint, groupID, memberID))
	if origin == "" {
		origin = "uaa"
	}
	if entityType == "" {
		entityType = "USER"
	}
	membership := GroupMember{Origin: origin, Type: entityType, Value: memberID}
	j, err := json.Marshal(membership)
	if err != nil {
		return err
	}
	err = a.doJSON(http.MethodDelete, &u, bytes.NewBuffer([]byte(j)), nil, true)
	if err != nil {
		return err
	}
	return nil
}

// GetGroupByName gets the group with the given name
// http://docs.cloudfoundry.org/api/uaa/version/4.14.0/index.html#list-4.
func (a *API) GetGroupByName(name string, attributes string) (*Group, error) {
	if name == "" {
		return nil, errors.New("group name may not be blank")
	}

	filter := fmt.Sprintf(`displayName eq "%v"`, name)
	groups, err := a.ListAllGroups(filter, "", attributes, "")
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("group %v not found", name)
	}
	return &groups[0], nil
}
