package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

func DecodeBytesToMapJson(bytes []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, errors.New("error decoding to json")
	}

	return result, nil
}

func GetStringValueOfKey(key string, data map[string]interface{}) (string, error) {
	value, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("error during getting value of %v", value)
	}
	return value, nil
}

func GetFloatValueOfKey(key string, data map[string]interface{}) (float64, error) {
	value, ok := data[key].(float64)
	if !ok {
		return 0.0, fmt.Errorf("error during getting value of %v", value)
	}
	return value, nil
}
