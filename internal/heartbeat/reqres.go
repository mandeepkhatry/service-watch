package heartbeat

import (
	"encoding/json"
	"fmt"
	"service-watch/internal/client"
	"service-watch/internal/generate"
	"service-watch/internal/utils"

	"github.com/getkin/kin-openapi/openapi3"
)

func SendHeartBeart(swaggerConfig *openapi3.Swagger, config map[string]interface{}) error {
	httpClient := client.NewHTTPClient(config)

	//App Specific Login
	loginResponse, _ := httpClient.ExecuteRequest("post", "/login", map[string]interface{}{"username": "enish_paneru", "password": "paneru"}, nil)

	access_token := "Bearer " + loginResponse.Message.(map[string]interface{})["data"].(map[string]interface{})["access_token"].(string)

	for ep, epProperties := range swaggerConfig.Paths {
		for method, methodProperties := range epProperties.Operations() {

			if methodProperties.RequestBody != nil {

				for content, contentProperties := range methodProperties.RequestBody.Value.Content {
					if content == "application/json" {
						component, subcomponent := utils.FindComponent(contentProperties.Schema.Ref)
						if component == "schemas" {
							var schema map[string]interface{}
							schemaBytes, _ := swaggerConfig.Components.Schemas[subcomponent].MarshalJSON()
							json.Unmarshal(schemaBytes, &schema)
							dummyData := generate.GenerateObject(schema)

							fmt.Println(methodProperties.Security)

							response, _ := httpClient.ExecuteRequest(method, ep, dummyData, map[string]string{"Authorization": access_token})

							fmt.Println("--------------------------------")
							fmt.Println("endpoint : ", ep)
							fmt.Println(response)
							fmt.Println("--------------------------------")

						}
					}
				}

			}

		}

	}

	return nil
}
