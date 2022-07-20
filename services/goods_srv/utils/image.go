package utils

import (
	"encoding/json"
)

func GetImageUrl(data []byte) []string {
	var res []string
	_ = json.Unmarshal(data, &res)
	return res
}

func SetImageUrl(urls []string) []byte {
	bytes, _ := json.Marshal(urls)
	return bytes
}
