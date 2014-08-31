package utils

import (
	"encoding/json"
)

func Marshal(v interface{}) (j string, err error) {
	bytes, err := json.Marshal(v)
	j = string(bytes)
	return
}
func UnMarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
