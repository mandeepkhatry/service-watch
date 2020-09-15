package utils

import "strings"

func DataMatching(data map[string]interface{}) map[string]interface{} {

	adjustedData := make(map[string]interface{})

	for k, v := range data {
		keySplits := strings.Split(k, "_")
		if keySplits[0] == "confirm" {
			adjustedData[k] = data[keySplits[1]]
		} else {
			adjustedData[k] = v
		}
	}

	return adjustedData
}
