package schema

import (
	"encoding/json"
	"service-watch/internal/generate"

	"github.com/getkin/kin-openapi/openapi3"
)

func GenerateSchemaData(schema *openapi3.SchemaRef, components openapi3.Components) map[string]interface{} {

	dummyData := make(map[string]interface{})

	if len(schema.Value.AnyOf) != 0 {
		for _, eachRef := range schema.Value.AnyOf {
			return GenerateSchemaData(eachRef, components)
		}
	} else if len(schema.Value.OneOf) != 0 {
		for _, eachRef := range schema.Value.OneOf {
			return GenerateSchemaData(eachRef, components)
		}

	} else if len(schema.Value.AllOf) != 0 {
		for _, eachRef := range schema.Value.AllOf {
			for k, v := range GenerateSchemaData(eachRef, components) {
				dummyData[k] = v
			}
		}
		return dummyData
	}

	innerSchema := make(map[string]interface{})

	schemaBytes, _ := schema.Value.MarshalJSON()

	json.Unmarshal(schemaBytes, &innerSchema)

	return generate.GenerateDummyData(innerSchema)

}
