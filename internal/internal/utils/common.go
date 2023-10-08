package utils

import (
	"encoding/json"

	"github.com/google/uuid"
)

func IsValidUrlId(urlId string) bool {
	if urlId == "" {
		return false
	}

	if _, err := uuid.Parse(urlId); err != nil {
		return false
	}

	return true
}

func ConvertMapToString(payload interface{}) (string, error) {
	stringPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(stringPayload), nil
}

func ConvertStringToMap(payload string) (map[string]interface{}, error) {
	var stringPayload map[string]interface{}
	err := json.Unmarshal([]byte(payload), &stringPayload)
	if err != nil {
		return nil, err
	}

	return stringPayload, nil
}
