package generate

//GenerateDummyData generates dummy data as per schema
func GenerateDummyData(schema map[string]interface{}) map[string]interface{} {
	return GenerateObject(schema)
}
