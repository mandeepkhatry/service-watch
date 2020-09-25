package heartbeat

import (
	"fmt"
	"service-watch/internal/client"
	"service-watch/internal/def"
	"service-watch/internal/models"
	"service-watch/internal/parser"
	"service-watch/internal/requestconfig"
	"service-watch/internal/schema"
	"service-watch/internal/security"
	"service-watch/internal/utils"
)

func ProcessRequest(appConfig models.AppConfig, config map[string]interface{}) error {

	httpClient := client.NewHTTPClient(config)

	endpointsDataBuffer := make(map[string]map[string]models.DataBuffer)

	securityScheme := security.HTTPAuthenticationScheme{HttpClient: httpClient}

	if security.ValidateSecuritySchemas(appConfig.Endpoints) {
		securityScheme.Credentials = def.HTTPSecurityCredentials
		securityScheme.Run()
		utils.DetachSecurityEndpoints(appConfig.Endpoints)
	}

	requestConfig := requestconfig.NewRequestConfig(securityScheme.Response)

	for _, childEps := range appConfig.Endpoints {

		for _, eachChildEp := range childEps {
			for childEpName, childEpProp := range eachChildEp {

				for _, eachMethod := range childEpProp.Methods {
					fmt.Println("--------------------------------")
					fmt.Println("Endpoint : ", childEpName)

					for methodName, methodOperations := range eachMethod {
						fmt.Println("Method : ", methodName)

						if _, present := def.SchemaBasedMethods[methodName]; present {

							if _, contentPresent := methodOperations.RequestBody.Value.Content["application/json"]; contentPresent {

								dummyData := schema.GenerateSchemaData(methodOperations.RequestBody.Value.Content["application/json"].Schema, appConfig.Components)

								dBuffer := models.DataBuffer{}

								dBuffer.AssignRequest(dummyData)

								specificEndpoint := parser.GenerateSpecificEndpoint(childEpName, endpointsDataBuffer, methodOperations.Parameters)

								response, _ := httpClient.ExecuteRequest(methodName, specificEndpoint, dummyData, requestConfig.Content)

								fmt.Println(response)

								dBuffer.AssignResponse(response.Message.(map[string]interface{}))

								if _, epNamePresent := endpointsDataBuffer[childEpName]; !epNamePresent {
									endpointsDataBuffer[childEpName] = make(map[string]models.DataBuffer)
								}

								endpointsDataBuffer[childEpName][methodName] = dBuffer

							}
						} else {

							//Methods GET/DELETE

							dBuffer := models.DataBuffer{}

							specificEndpoint := parser.GenerateSpecificEndpoint(childEpName, endpointsDataBuffer, methodOperations.Parameters)

							response, _ := httpClient.ExecuteRequest(methodName, specificEndpoint, nil, requestConfig.Content)

							fmt.Println("----specific---- : ", specificEndpoint)
							//Intercept for query
							if len(methodOperations.Parameters) != 0 {
								ep := parser.GenerateRequestQuery(specificEndpoint, methodOperations.Parameters, response)

								response, _ = httpClient.ExecuteRequest(methodName, ep, nil, requestConfig.Content)

							}
							fmt.Println(response)

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

	return nil

}
