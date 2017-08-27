package uaa

import (
	"encoding/json"
	"net/http"
)

type Userinfo struct {
	UserId string `json:"user_id"`
	Sub string
	Username string `json:"user_name"`
	GivenName string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email string
	PhoneNumber []string
	PreviousLoginTime int64
	Name string
}

func Me(client *http.Client, config Config) (Userinfo, error) {
	body, err := AuthenticatedRequester{}.GetBytes(client, config, "/userinfo", "scheme=openid")
	if err != nil {
		return Userinfo{}, err
	}

	info := Userinfo{}
	err = json.Unmarshal(body,&info)
	if err != nil {
		return Userinfo{}, parseError("/userinfo", body)
	}

	return info, nil
}