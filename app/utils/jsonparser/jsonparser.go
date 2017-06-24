package jsonparser

import (
	"encoding/json"
	"fmt"
	"os"
)

func Parse(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func ToJsonString(c interface{}) (string, error) {
	bytes, err := json.Marshal(c)
	return string(bytes), err
}
