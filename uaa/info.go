package uaa

import (
	"net/http"
	"github.com/jhamon/guac/utils"

	"io/ioutil"
	"encoding/json"
	"errors"
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

func Info(context UaaContext) (UaaInfo, error) {
	infoUrl := utils.BuildUrl(context.BaseUrl, "info").String()

	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", infoUrl, nil)
	req.Header.Add("Accept","application/json")

	resp, err := httpClient.Do(req)
	if (resp.StatusCode != 200 || err != nil) {
		return UaaInfo{}, requestError(infoUrl)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UaaInfo{}, unknownError()
	}

	info := UaaInfo{}
	err = json.Unmarshal(body,&info)
	if err != nil {
		return UaaInfo{}, parseError(infoUrl, body)
	}

	return info, nil
}

func requestError(url string) error {
	return errors.New("An unknown error occurred while calling " + url)
}

func parseError(url string, body []byte) error {
	errorMsg := "An unknown error occurred while parsing response from " + url + ". Response was " + string(body)
	return errors.New(errorMsg)
}

func unknownError() error {
	return errors.New("An unknown error occurred")
}