package uaa

import (
	"encoding/json"
	"net/http"
)

type UaaInfo struct {
	App            uaaApp              `json:"app"`
	Links          uaaLinks            `json:"links"`
	Prompts        map[string][]string `json:"prompts"`
	ZoneName       string              `json:"zone_name"`
	EntityId       string              `json:"entityID"`
	CommitId       string              `json:"commit_id"`
	Timestamp      string              `json:"timestamp"`
	IdpDefinitions map[string]string   `json:"idpDefinitions"`
}

type uaaApp struct {
	Version string `json:"version"`
}

type uaaLinks struct {
	ForgotPassword string `json:"passwd"`
	Uaa            string `json:"uaa"`
	Registration   string `json:"register"`
	Login          string `json:"login"`
}

func Info(client *http.Client, config Config) (UaaInfo, error) {
	bytes, err := UnauthenticatedRequester{}.Get(client, config, "info", "")
	if err != nil {
		return UaaInfo{}, err
	}

	info := UaaInfo{}
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return UaaInfo{}, parseError("", bytes)
	}

	return info, err
}
