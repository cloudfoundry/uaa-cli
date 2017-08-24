package uaa

import (
	"encoding/json"
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

func Me(context UaaContext) (Userinfo, error) {
	body, err := AuthenticatedGetter{}.Get(context, "/userinfo", "scheme=openid")
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