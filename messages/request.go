package messages

import (
	"errors"
	"fmt"

	"encoding/json"
	"io"
	"io/ioutil"
)

func RequestBodyToMap(requestBody io.ReadCloser) (map[string]interface{}, error) {
	if requestBody == nil {
		return nil, errors.New("Empty body")
	}

	rawData, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return nil, err
	}

	params := make(map[string]interface{})
	err = json.Unmarshal(rawData, &params)
	fmt.Printf(">>> body data: %v\n", params)

	return params, err
}
