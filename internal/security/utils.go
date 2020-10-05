package security

import (
	"service-watch/internal/models"
)

func ValidateSecuritySchemas(endpoints map[string][]map[string]models.Endpoint, securityEndpoints []string, credentials map[string]interface{}) bool {

	for _, securityEp := range securityEndpoints {
		if _, present := endpoints[securityEp]; !present {
			return false
		}
	}

	return true
}
