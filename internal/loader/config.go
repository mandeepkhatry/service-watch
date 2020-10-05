package loader

import (
	"service-watch/internal/def"
	"service-watch/internal/models"
	"service-watch/internal/parser"
	"service-watch/internal/utils"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadSwagger(swaggerConfig *openapi3.Swagger, watchConfig map[string]interface{}) (*models.AppConfig, error) {

	appConfig := &models.AppConfig{}

	//Form : RootEndpoint : array of childs
	endpoints := make(map[string][]map[string]models.Endpoint)

	ignoreEndpoints := watchConfig["ignore_endpoints"].([]interface{})

	endpointsIngnored := make(map[string]bool)

	for _, v := range ignoreEndpoints {
		endpointsIngnored[v.(string)] = true
	}

	for ep, _ := range swaggerConfig.Paths {

		if _, present := endpointsIngnored[ep]; !present {
			endpointParts := parser.ParseEndpoint(ep)

			root := endpointParts[0]

			if _, present := endpoints[root]; !present {
				endpoints[root] = []map[string]models.Endpoint{}
			}

			for _, eachEpPart := range endpointParts {

				if !utils.EpRedundancyPresent(eachEpPart, endpoints[root]) {
					endpoints[root] = append(endpoints[root], map[string]models.Endpoint{eachEpPart: models.Endpoint{Root: eachEpPart}})

				}

			}
		}

	}

	appSpecificEndpoints := make(map[string][]map[string]models.Endpoint)

	for root, eachChildEndpoints := range endpoints {

		appSpecificEndpoints[root] = []map[string]models.Endpoint{}

		bufferEndpoints := make([]map[string]models.Endpoint, 0)

		for _, eachChildEp := range eachChildEndpoints {

			bufferChildEp := make(map[string]models.Endpoint)

			for eachChildEpName, eachChildProp := range eachChildEp {
				if _, epPresent := swaggerConfig.Paths[eachChildEpName]; epPresent {

					bufferModel := &models.Endpoint{}

					for _, eachOrderedMethod := range def.OrderedMethods {

						if _, operationPresent := swaggerConfig.Paths[eachChildEpName].Operations()[eachOrderedMethod]; operationPresent {
							bufferModel.Methods = append(bufferModel.Methods, map[string]*openapi3.Operation{eachOrderedMethod: swaggerConfig.Paths[eachChildEpName].Operations()[eachOrderedMethod]})
						}
					}

					bufferModel.Root = eachChildProp.Root

					bufferChildEp[eachChildEpName] = *bufferModel

				}
			}

			bufferEndpoints = append(bufferEndpoints, bufferChildEp)

		}

		appSpecificEndpoints[root] = bufferEndpoints

	}

	appConfig.Endpoints = appSpecificEndpoints
	appConfig.Server = swaggerConfig.Servers[0].URL
	appConfig.Components = swaggerConfig.Components

	return appConfig, nil

}
