package fixtures

import "encoding/json"

func EntityResponse(response interface{}) string {
	bytes, err := json.Marshal(response)
	if err != nil {
		return "\"Failed to marshal provided entity\""
	}

	return string(bytes)
}
