package utils

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/araddon/dateparse"
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

func ConvertBuffer(buffer *bytes.Buffer) map[string]interface{} {
	data := make(map[string]interface{})
	json.Unmarshal(buffer.Bytes(), &data)
	return data
}

func ConvertMap(data map[string]interface{}) *bytes.Buffer {
	dataBytes, _ := json.Marshal(data)
	return bytes.NewBuffer(dataBytes)
}

func FormatConstantDate(s string) (string, error) {
	var err error
	if dateFormat, err := dateparse.ParseFormat(s); err == nil {
		return dateFormat, nil
	}
	return s, err
}
