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

func ProcessRequest(appConfig models.AppConfig, config map[string]interface{}) (map[string][]interface{}, error) {

	httpClient := client.NewHTTPClient(config)

	endpointsDataBuffer := make(map[string]map[string]models.DataBuffer)

	securityScheme := security.HTTPAuthenticationScheme{HttpClient: httpClient}

	if security.ValidateSecuritySchemas(appConfig.Endpoints) {
		securityScheme.Credentials = def.HTTPSecurityCredentials
		securityScheme.Run()
		utils.DetachSecurityEndpoints(appConfig.Endpoints)
	}

	/**
	Logs format
		key : status ["timeout", "success", "failure"]
		value : array of corresponding response

	**/
	logs := make(map[string][]interface{})

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

							contentType := utils.GetContent(methodOperations.RequestBody.Value.Content)

							buffer, contentType, _ := content.ContentBasedData[contentType](methodOperations.RequestBody.Value.Content[contentType].Schema, methodOperations.RequestBody.Value.Content[contentType].Encoding, appConfig.Components)

							dBuffer := models.DataBuffer{}

							dBuffer.AssignRequest(utils.ConvertBuffer(buffer))

							specificEndpoint := parser.GenerateSpecificEndpoint(childEpName, endpointsDataBuffer, methodOperations.Parameters)

							requestConfig.Content["Content-Type"] = contentType

							response, _ := httpClient.ExecuteRequest(methodName, specificEndpoint, buffer, requestConfig.Content)

							//Find status of endpoint and append to logs accordingly

							status, log := utils.LogExtract(childEpName, methodName, response)

							if _, logKeyPresent := logs[status]; !logKeyPresent {
								logs[status] = make([]interface{}, 0)
							}

							logs[status] = append(logs[status], log)

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

								//Find status of endpoint and append to logs accordingly

								status, log := utils.LogExtract(childEpName, methodName, response)

								if _, logKeyPresent := logs[status]; !logKeyPresent {
									logs[status] = make([]interface{}, 0)
								}

								logs[status] = append(logs[status], log)

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
		fmt.Println(response)
		//Find status of endpoint and append to logs accordingly

		status, log := utils.LogExtract(ep, epProperties.MethodName, response)

		if _, logKeyPresent := logs[status]; !logKeyPresent {
			logs[status] = make([]interface{}, 0)
		}

		logs[status] = append(logs[status], log)
	}

	return logs, nil

}
