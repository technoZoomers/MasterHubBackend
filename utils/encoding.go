package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/logger"
)

func Encode(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		parseError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Errorf(parseError.Error())
		return "", parseError
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func Decode(in string, obj interface{}) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		parseError := fmt.Errorf("error decoding string: %v", err.Error())
		logger.Errorf(parseError.Error())
		return parseError
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		parseError := fmt.Errorf("error unmarshalling json: %v", err.Error())
		logger.Errorf(parseError.Error())
		return parseError
	}
	return nil
}
