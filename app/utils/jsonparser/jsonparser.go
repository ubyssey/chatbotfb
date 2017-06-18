package jsonparser

import (
	"encoding/json"
	"fmt"
	"os"
)

func Parse(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func ToJsonString(c interface{}) string {
	bytes, err := json.Marshal(c)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}
