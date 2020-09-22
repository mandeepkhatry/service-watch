package utils

import (
	"strings"
)

func DataMatching(data map[string]interface{}) map[string]interface{} {

	adjustedData := make(map[string]interface{})

	for k, v := range data {
		keySplits := strings.Split(k, "_")
		if keySplits[0] == "confirm" {
			adjustedData[k] = data[strings.Join(keySplits[1:], "_")]
		} else {
			adjustedData[k] = v
		}
	}

	return adjustedData
}
