package utils

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
)

func GetCorrespondingSchema(status string, responses openapi3.Responses, components openapi3.Components) map[string]interface{} {

	contentType := GetContent(responses[status].Value.Content)

	openapiSchema := responses[status].Value.Content[contentType].Schema

	component, subcomponent := FindComponent(openapiSchema.Ref)

	var schema map[string]interface{}

	if component == "schemas" {
		schemaBytes, _ := components.Schemas[subcomponent].MarshalJSON()
		json.Unmarshal(schemaBytes, &schema)

	}

	return schema

}
