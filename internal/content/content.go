package content

import (
	"bytes"
	"encoding/json"
	"service-watch/internal/schema"

	"github.com/getkin/kin-openapi/openapi3"
)

var ContentBasedData = map[string]func(configSchema *openapi3.SchemaRef, encoding map[string]*openapi3.Encoding, components openapi3.Components) (*bytes.Buffer, error){
	"application/json": func(configSchema *openapi3.SchemaRef, encoding map[string]*openapi3.Encoding, components openapi3.Components) (*bytes.Buffer, error) {
		dummyData := schema.GenerateSchemaData(configSchema, components)
		requestBytes, err := json.Marshal(dummyData)
		if err != nil {
			return bytes.NewBuffer([]byte{}), err
		}
		return bytes.NewBuffer(requestBytes), nil
	},
	"multipart/form-data": func(configSchema *openapi3.SchemaRef, encoding map[string]*openapi3.Encoding, components openapi3.Components) (*bytes.Buffer, error) {
		dummyData := schema.GenerateSchemaData(configSchema, components)

		requestBytes, err := json.Marshal(dummyData)
		if err != nil {
			return bytes.NewBuffer([]byte{}), err
		}
		return bytes.NewBuffer(requestBytes), nil
	},
}
