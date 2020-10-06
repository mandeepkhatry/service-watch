package utils

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
)

func GetCorrespondingSchema(status string, responses openapi3.Responses) map[string]interface{} {

	responseByte, _ := responses[status].Value.MarshalJSON()

	var responseMap map[string]interface{}

	json.Unmarshal(responseByte, &responseMap)

	contentType := ""

	for cType := range responseMap["content"].(map[string]interface{}) {
		contentType = cType
		break
	}

	openapiSchema := responses[status].Value.Content[contentType].Schema

	var schema map[string]interface{}

	schemaByte, _ := openapiSchema.Value.MarshalJSON()

	json.Unmarshal(schemaByte, &schema)

	return schema

}
