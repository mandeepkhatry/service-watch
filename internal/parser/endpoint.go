package parser

import (
	"encoding/json"
	"fmt"
	"service-watch/internal/generate"
	"service-watch/internal/models"
	"service-watch/internal/utils"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func ParseEndpoint(endpoint string) []string {

	epParts := strings.Split(endpoint, "/")

	roots := make([]string, 0)

	if epParts[0] == "" {
		epParts = epParts[1:]
	}

	for i := range epParts {

		roots = append(roots, "/"+strings.Join(epParts[0:i+1], "/"))
	}

	return roots
}

func GenerateSpecificEndpoint(ep string, endpointsDataBuffer map[string]map[string]models.DataBuffer, parameters openapi3.Parameters) string {

	epParts := strings.Split(ep, "/")

	finalEp := "/"

	oneTimePost := true

	current := 0

	epParts = epParts[1:]

	for i, epPart := range epParts {
		if utils.ValidateResource(epPart) {

			root := finalEp + strings.Join(epParts[current:i], "/")

			resource := ""

			if oneTimePost {
				//epPart is resource since already validated
				resource = endpointsDataBuffer[root]["POST"].Response["data"].(map[string]interface{})[utils.SeperateResource(epPart)].(string)
				oneTimePost = false
			} else {
				resource = endpointsDataBuffer[root]["GET"].Response["data"].(map[string]interface{})[utils.SeperateResource(epPart)].(string)
			}

			finalEp = root + "/" + resource + "/"

			if i <= len(epPart)-1 {
				current = i + 1
			}

		}

	}

	if finalEp == "/" {
		finalEp = ep
	} else {
		if !utils.ValidateResource(epParts[len(epParts)-1]) {
			finalEp += epParts[len(epParts)-1]
		} else {
			finalEp = strings.TrimSuffix(finalEp, "/")
		}
	}

	return finalEp

}

func GenerateRequestQuery(ep string, parameters openapi3.Parameters, response models.HeartBeatResponse) string {

	//Take initial response object
	res := response.Message.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})

	querySegments := make([]string, 0)

	for _, params := range parameters {

		if params.Value.Schema != nil {
			//Take data according to schema
			var schema map[string]interface{}
			schemaBytes, _ := params.Value.Schema.MarshalJSON()
			json.Unmarshal(schemaBytes, &schema)

			generatedData := generate.GenerateDummyData(schema)

			querySegments = append(querySegments, params.Value.Name+"="+fmt.Sprintf("%v", generatedData[params.Value.Name]))
		} else {
			//Take data from root endpoint
			field := params.Value.Name
			querySegments = append(querySegments, field+"="+fmt.Sprintf("%v", res[field]))
		}

	}

	query := strings.Join(querySegments, "&")

	return ep + "?" + query

}
