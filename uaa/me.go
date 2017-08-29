package uaa

import (
	"encoding/json"
	"net/http"
)

type Userinfo struct {
	UserId            string   `json:"user_id"`
	Sub               string   `json:"sub"`
	Username          string   `json:"user_name"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	Email             string   `json:"email"`
	PhoneNumber       []string `json:"phone_number"`
	PreviousLoginTime int64    `json:"previous_logon_time"`
	Name              string   `json:"name"`
}

func Me(client *http.Client, config Config) (Userinfo, error) {
	body, err := AuthenticatedRequester{}.Get(client, config, "/userinfo", "scheme=openid")
	if err != nil {
		return Userinfo{}, err
	}

	info := Userinfo{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return Userinfo{}, parseError("/userinfo", body)
	}

	return info, nil
}
