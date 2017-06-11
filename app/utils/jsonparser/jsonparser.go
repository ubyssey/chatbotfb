package jsonparser

import (
	"encoding/json"
	"error"
	"string"
)

// Parse to JSON
func (c *interface{}) Parse(b []byte) error {
	return json.Unmarshal(b, &c)
}

func ToJson(c interface{}) string {
	bytes, err := json.Marshal(c)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}
