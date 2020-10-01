package utils

import (
	"encoding/json"
	"service-watch/internal/def"
	"service-watch/internal/mime"

	"github.com/getkin/kin-openapi/openapi3"
)

/*
FindFileContent returns map[string]string

	Format :
		key : field
		value : file_type

*/
func FindFileContent(configSchema *openapi3.SchemaRef, components openapi3.Components, encoding map[string]*openapi3.Encoding) map[string]string {

	fileContent := make(map[string]string)

	component, subcomponent := FindComponent((configSchema.Ref))
	if component == "schemas" {
		var schema map[string]interface{}
		schemaBytes, _ := components.Schemas[subcomponent].MarshalJSON()
		json.Unmarshal(schemaBytes, &schema)

		if _, present := schema["properties"]; present {
			objectProperties := schema["properties"].(map[string]interface{})

			for field, fieldProperties := range objectProperties {

				if _, formatPresent := fieldProperties.(map[string]interface{})["format"].(string); formatPresent {
					format := fieldProperties.(map[string]interface{})["format"].(string)

					if _, encodingPresent := def.ContentEncoding[format]; encodingPresent {
						if _, openAPIEncoding := encoding[field]; openAPIEncoding {
							fileContent[field] = mime.GetExtension(encoding[field].ContentType)
						}
					}
				}

			}
		}
	}
	return fileContent

}
