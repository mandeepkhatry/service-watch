package schema

import (
	"encoding/json"
	"service-watch/internal/generate"
	"service-watch/internal/utils"

	"github.com/getkin/kin-openapi/openapi3"
)

func GenerateSchemaData(schema *openapi3.SchemaRef, components openapi3.Components) map[string]interface{} {

	dummyData := make(map[string]interface{})

	if len(schema.Ref) != 0 {
		component, subcomponent := utils.FindComponent((schema.Ref))

		if component == "schemas" {
			var schema map[string]interface{}
			schemaBytes, _ := components.Schemas[subcomponent].MarshalJSON()

			json.Unmarshal(schemaBytes, &schema)

			return generate.GenerateDummyData(schema)

		}
	}

	for _, eachRef := range schema.Value.OneOf {
		return GenerateSchemaData(eachRef, components)
	}

	for _, eachRef := range schema.Value.AnyOf {
		return GenerateSchemaData(eachRef, components)
	}

	for _, eachRef := range schema.Value.AllOf {
		for k, v := range GenerateSchemaData(eachRef, components) {
			dummyData[k] = v
		}
	}

	return dummyData
}
