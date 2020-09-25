package heartbeat

import (
	"fmt"
	"service-watch/internal/client"
	"service-watch/internal/content"
	"service-watch/internal/def"
	"service-watch/internal/models"
	"service-watch/internal/parser"
	"service-watch/internal/requestconfig"
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

	//Deletion Methods
	rearMostEndpointCollection := make(map[string]models.RearMostEndpoint)

	requestConfig := requestconfig.NewRequestConfig(securityScheme.Response)

	for _, childEps := range appConfig.Endpoints {

		for _, eachChildEp := range childEps {
			for childEpName, childEpProp := range eachChildEp {
				fmt.Println("--------------------------------")
				fmt.Println("Endpoint : ", childEpName)
				for _, eachMethod := range childEpProp.Methods {

					for methodName, methodOperations := range eachMethod {
						fmt.Println("Method : ", methodName)

						if _, present := def.SchemaBasedMethods[methodName]; present {

							contentType := content.GetContent(methodOperations.RequestBody.Value.Content)

							buffer, _ := content.ContentBasedData[contentType](methodOperations.RequestBody.Value.Content[contentType].Schema, appConfig.Components)

							dBuffer := models.DataBuffer{}

							dBuffer.AssignRequest(utils.ConvertBuffer(buffer))

							specificEndpoint := parser.GenerateSpecificEndpoint(childEpName, endpointsDataBuffer, methodOperations.Parameters)

							response, _ := httpClient.ExecuteRequest(methodName, specificEndpoint, buffer, requestConfig.Content)

							fmt.Println(response)

							dBuffer.AssignResponse(response.Message.(map[string]interface{}))

							if _, epNamePresent := endpointsDataBuffer[childEpName]; !epNamePresent {
								endpointsDataBuffer[childEpName] = make(map[string]models.DataBuffer)
							}

							endpointsDataBuffer[childEpName][methodName] = dBuffer

						} else {

							//Methods GET/DELETE

							if _, present := def.RearMostMethod[methodName]; present {

								specificEndpoint := parser.GenerateSpecificEndpoint(childEpName, endpointsDataBuffer, methodOperations.Parameters)
								rearMostEp := models.RearMostEndpoint{
									MethodName:       methodName,
									SpecificEndpoint: specificEndpoint,
									RequestBody:      nil,
									RequestConfig:    requestConfig.Content,
								}

								rearMostEndpointCollection[childEpName] = rearMostEp

							} else {

								dBuffer := models.DataBuffer{}

								specificEndpoint := parser.GenerateSpecificEndpoint(childEpName, endpointsDataBuffer, methodOperations.Parameters)

								response, _ := httpClient.ExecuteRequest(methodName, specificEndpoint, nil, requestConfig.Content)

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

	}

	//Deletion Method considered as RearMostEnd
	for ep, epProperties := range rearMostEndpointCollection {
		fmt.Println("--------------------------------")
		fmt.Println("Endpoint : ", ep)
		fmt.Println("Method : ", epProperties.MethodName)
		response, _ := httpClient.ExecuteRequest(epProperties.MethodName, epProperties.SpecificEndpoint, utils.ConvertMap(epProperties.RequestBody), epProperties.RequestConfig)
		fmt.Println("----specific---- : ", epProperties.SpecificEndpoint)
		fmt.Println(response)
	}

	return nil

}
