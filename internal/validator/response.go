package validator

import (
	"github.com/xeipuuv/gojsonschema"
)

func ValidateResponseWrtSchema(schema map[string]interface{}, data map[string]interface{}) (bool, map[string]interface{}, error) {
	loader := gojsonschema.NewGoLoader(schema)
	dataLoader := gojsonschema.NewGoLoader(data)

	errorMap := make(map[string]interface{})

	result, err := gojsonschema.Validate(loader, dataLoader)
	if err != nil {
		return false, errorMap, err
	}

	if result.Valid() {
		return true, errorMap, nil
	}

	for _, v := range result.Errors() {
		errorMap[v.Field()] = v.Description()
	}

	return false, errorMap, nil

}
