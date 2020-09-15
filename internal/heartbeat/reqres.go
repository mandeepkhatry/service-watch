package heartbeat

import (
	"encoding/json"
	"fmt"
	"service-watch/internal/client"
	"service-watch/internal/def"
	"service-watch/internal/generate"
	"service-watch/internal/models"
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
							fmt.Println(childEpName, "-----", methodName)
							if _, present := def.SchemaBasedMethods[methodName]; present {
								if _, contentPresent := methodOperations.RequestBody.Value.Content["application/json"]; contentPresent {

									component, subcomponent := utils.FindComponent(methodOperations.RequestBody.Value.Content["application/json"].Schema.Ref)

									if component == "schemas" {

										var schema map[string]interface{}

										schemaBytes, _ := appConfig.Components.Schemas[subcomponent].MarshalJSON()

										json.Unmarshal(schemaBytes, &schema)
										dummyData := generate.GenerateObject(schema)

										dBuffer := models.DataBuffer{}

										dBuffer.AssignRequest(dummyData)

										//TODOImplment Security Options. Not present now in openapi test config. Also Protocol on Reponse Level has to be confirmed.

										rootEp, resource, isSpecific := utils.IsSpecificItem(childEpName)

										if isSpecific {

											//Check for root response for resource ids.

											resourceItem := endpointsDataBuffer[rootEp]["POST"].Response["data"].(map[string]interface{})[resource].(string)

											resourceDisplacedEp := rootEp + "/" + resourceItem

											response, _ := httpClient.ExecuteRequest(methodName, resourceDisplacedEp, dummyData, map[string]string{"Authorization": access_token})

											dBuffer.AssignResponse(response.Message.(map[string]interface{}))

											endpointsDataBuffer[childEpName] = make(map[string]models.DataBuffer)

											endpointsDataBuffer[childEpName][methodName] = dBuffer

										} else {
											response, _ := httpClient.ExecuteRequest(methodName, childEpName, dummyData, map[string]string{"Authorization": access_token})

											dBuffer.AssignResponse(response.Message.(map[string]interface{}))

											endpointsDataBuffer[childEpName] = make(map[string]models.DataBuffer)

											endpointsDataBuffer[childEpName][methodName] = dBuffer

										}

										fmt.Println(dBuffer)

									}

								}
							} else {
								//Methods GET/DELETE
								dBuffer := models.DataBuffer{}
								rootEp, resource, isSpecific := utils.IsSpecificItem(childEpName)
								if isSpecific {

									//Check for root response for resource ids.

									resourceItem := endpointsDataBuffer[rootEp]["POST"].Response["data"].(map[string]interface{})[resource].(string)

									resourceDisplacedEp := rootEp + "/" + resourceItem

									fmt.Println(resourceDisplacedEp)

									response, _ := httpClient.ExecuteRequest(methodName, resourceDisplacedEp, nil, map[string]string{"Authorization": access_token})

									dBuffer.AssignResponse(response.Message.(map[string]interface{}))

									endpointsDataBuffer[childEpName] = make(map[string]models.DataBuffer)
									endpointsDataBuffer[childEpName][methodName] = dBuffer

								} else {
									response, _ := httpClient.ExecuteRequest(methodName, childEpName, nil, map[string]string{"Authorization": access_token})

									dBuffer.AssignResponse(response.Message.(map[string]interface{}))

									endpointsDataBuffer[childEpName] = make(map[string]models.DataBuffer)
									endpointsDataBuffer[childEpName][methodName] = dBuffer

								}
								fmt.Println(dBuffer)

							}

						}

					}
				}
			}

		}

	}

	return nil

}
