package info

import (

)
import (
	"net/http"
	"github.com/jhamon/guac/utils"
	"io/ioutil"
	"encoding/json"
)

type JWK struct {
	Kty string
	E string
	Use string
	Kid string
	Alg string
	Value string
	N string
}

func TokenKey(context UaaContext) (JWK, error) {
	tokenKeyUrl := utils.BuildUrl(context.BaseUrl, "token_key")
	url := tokenKeyUrl.String()

	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept","application/json")

	resp, err := httpClient.Do(req)
	if (resp.StatusCode != 200 || err != nil) {
		return JWK{}, requestError(url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return JWK{}, unknownError()
	}

	key := JWK{}
	err = json.Unmarshal(body,&key)
	if err != nil {
		return JWK{}, parseError(url, body)
	}

	return key, nil
}