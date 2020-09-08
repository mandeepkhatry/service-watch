package generate

//FieldToGenerator maps field type to corresponding field generator
var FieldToGenerator = map[string]func(properties map[string]interface{}) interface{}{
	"integer": func(properties map[string]interface{}) interface{} {
		return GenerateInteger(properties)
	},
	"number": func(properties map[string]interface{}) interface{} {
		return GenerateFloat(properties)
	},
	"string": func(properties map[string]interface{}) interface{} {
		return GenerateString(properties)
	},
	"array_string": func(properties map[string]interface{}) interface{} {
		return GenerateStringArray(properties)
	},
	"array_number": func(properties map[string]interface{}) interface{} {
		return GenerateNumberArray(properties)
	},
}
