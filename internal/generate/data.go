package generate

import "service-watch/internal/utils"

//GenerateDummyData generates dummy data as per schema
func GenerateDummyData(schema map[string]interface{}) map[string]interface{} {

	return utils.DataMatching(GenerateObject(schema))
}
