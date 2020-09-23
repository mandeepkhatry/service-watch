package security

import (
	"service-watch/internal/def"
	"service-watch/internal/models"
)

func ValidateSecuritySchemas(endpoints map[string][]map[string]models.Endpoint) bool {

	for _, securityEp := range def.HTTPSecurityEndpoints {
		if _, present := endpoints[securityEp]; !present {
			return false
		}
	}

	return true
}
