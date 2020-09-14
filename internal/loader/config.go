package loader

import (
	"service-watch/internal/def"
	"service-watch/internal/models"
	"service-watch/internal/parser"
	"service-watch/internal/utils"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadSwagger(swaggerConfig *openapi3.Swagger) (*models.AppConfig, error) {

	appConfig := &models.AppConfig{}

	//Form : RootEndpoint : array of childs
	endpoints := make(map[string][]map[string]models.Endpoint)

	for ep, _ := range swaggerConfig.Paths {

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

	for _, eachChildEndpoints := range endpoints {

		for _, eachChildEp := range eachChildEndpoints {

			for eachChildEpName, eachChildProp := range eachChildEp {
				if _, epPresent := swaggerConfig.Paths[eachChildEpName]; epPresent {

					for _, eachOrderedMethod := range def.OrderedMethods {

						if _, operationPresent := swaggerConfig.Paths[eachChildEpName].Operations()[eachOrderedMethod]; operationPresent {
							eachChildProp.Methods = append(eachChildProp.Methods, map[string]*openapi3.Operation{eachOrderedMethod: swaggerConfig.Paths[eachChildEpName].Operations()[eachOrderedMethod]})
						}
					}
					eachChildEp[eachChildEpName] = eachChildProp

				}
			}
		}

	}

	appConfig.Endpoints = endpoints
	appConfig.Server = swaggerConfig.Servers[0].URL

	return appConfig, nil

}
