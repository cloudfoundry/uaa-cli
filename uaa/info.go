package uaa

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
)

type UaaInfo struct {
	App uaaApp `json:"app"`
	Links uaaLinks `json:"links"`
	Prompts uaaPrompts `json:"prompts"`
	ZoneName string `json:"zone_name"`
	EntityId string `json:"entityID"`
	CommitId string `json:"commit_id"`
	Timestamp string `json:"timestamp"`
}

type uaaApp struct {
	Version string `json:"version"`
}

type uaaLinks struct {
	ForgotPassword string `json:"passwd"`
	Uaa string `json:"uaa"`
	Registration string `json:"register"`
	Login string `json:"login"`
}

type uaaPrompts struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
}

func Info(context UaaContext, client *http.Client) (UaaInfo, error) {
	req, err := UnauthenticatedRequestFactory{}.Get(context, "info", "")
	if err != nil {
		return UaaInfo{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return UaaInfo{}, requestError(req.URL.String())
	}

	if resp.StatusCode != http.StatusOK {
		return UaaInfo{}, requestError(req.URL.String())
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UaaInfo{}, err
	}

	info := UaaInfo{}
	err = json.Unmarshal(bytes,&info)
	if err != nil {
		return UaaInfo{}, parseError(req.URL.String(), bytes)
	}

	return info, err
}