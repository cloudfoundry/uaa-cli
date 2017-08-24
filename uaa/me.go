package uaa

import (
	"io/ioutil"
	"net/http"
	"github.com/jhamon/uaa-cli/utils"
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
	userinfoUrl, err := utils.BuildUrl(context.BaseUrl, "userinfo")
	if err != nil {
		return Userinfo{}, err
	}
	userinfoUrl.RawQuery = "scheme=openid"
	url := userinfoUrl.String()

	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept","application/json")
	req.Header.Add("Authorization","bearer " + context.AccessToken)

	resp, err := httpClient.Do(req)
	if (resp.StatusCode != 200 || err != nil) {
		return Userinfo{}, requestError(url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Userinfo{}, unknownError()
	}

	info := Userinfo{}
	err = json.Unmarshal(body,&info)
	if err != nil {
		return Userinfo{}, parseError(url, body)
	}

	return info, nil
}