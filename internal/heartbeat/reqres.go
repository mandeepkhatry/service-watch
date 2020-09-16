package heartbeat

import (
	"encoding/json"
	"service-watch/internal/client"
	"service-watch/internal/def"
	"service-watch/internal/generate"
	"service-watch/internal/models"
	"service-watch/internal/parser"
	"service-watch/internal/utils"
)

func ProcessRequest(appConfig models.AppConfig, config map[string]interface{}) error {
	httpClient := client.NewHTTPClient(config)

	endpointsDataBuffer := make(map[string]map[string]models.DataBuffer)

	//App Specific Login
	loginResponse, _ := httpClient.ExecuteRequest("POST", "/login", map[string]interface{}{"username": "enish_paneru", "password": "paneru"}, nil)

	access_token := "Bearer " + loginResponse.Message.(map[string]interface{})["data"].(map[string]interface{})["access_token"].(string)

	for root, childEps := range appConfig.Endpoints {

		if root != "/login" {
			for _, eachChildEp := range childEps {
				for childEpName, childEpProp := range eachChildEp {

					for _, eachMethod := range childEpProp.Methods {

						for methodName, methodOperations := range eachMethod {

							if _, present := def.SchemaBasedMethods[methodName]; present {
								if _, contentPresent := methodOperations.RequestBody.Value.Content["application/json"]; contentPresent {

									component, subcomponent := utils.FindComponent(methodOperations.RequestBody.Value.Content["application/json"].Schema.Ref)

									if component == "schemas" {

										var schema map[string]interface{}

										schemaBytes, _ := appConfig.Components.Schemas[subcomponent].MarshalJSON()

										json.Unmarshal(schemaBytes, &schema)
										dummyData := generate.GenerateDummyData(schema)

										dBuffer := models.DataBuffer{}

										dBuffer.AssignRequest(dummyData)

										specificEndpoint := parser.GenerateSpecificEndpoint(childEpName, endpointsDataBuffer, methodOperations.Parameters)

										response, _ := httpClient.ExecuteRequest(methodName, specificEndpoint, dummyData, map[string]string{"Authorization": access_token})

										dBuffer.AssignResponse(response.Message.(map[string]interface{}))

										if _, epNamePresent := endpointsDataBuffer[childEpName]; !epNamePresent {
											endpointsDataBuffer[childEpName] = make(map[string]models.DataBuffer)
										}

										endpointsDataBuffer[childEpName][methodName] = dBuffer

									}

								}
							} else {
								//Methods GET/DELETE

								dBuffer := models.DataBuffer{}

								specificEndpoint := parser.GenerateSpecificEndpoint(childEpName, endpointsDataBuffer, methodOperations.Parameters)

								response, _ := httpClient.ExecuteRequest(methodName, specificEndpoint, nil, map[string]string{"Authorization": access_token})

								dBuffer.AssignResponse(response.Message.(map[string]interface{}))

								if _, epNamePresent := endpointsDataBuffer[childEpName]; !epNamePresent {
									endpointsDataBuffer[childEpName] = make(map[string]models.DataBuffer)
								}

								endpointsDataBuffer[childEpName][methodName] = dBuffer

							}

						}

					}

				}
			}

		}

	}

	return nil

}
