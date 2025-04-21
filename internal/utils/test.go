package utils

import (
	"bytes"
	"encoding/json"
)

func MapToBuffer(m map[string]string) *bytes.Buffer {
	json, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(json)
}
