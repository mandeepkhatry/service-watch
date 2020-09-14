package utils

import (
	"service-watch/internal/models"
	"strings"
)

func ValidateResource(epPart string) bool {
	if strings.HasPrefix(epPart, "{") && strings.HasSuffix(epPart, "}") {
		return true
	}
	return false
}

func EpRedundancyPresent(epoint string, epoints []map[string]models.Endpoint) bool {
	for _, eachEp := range epoints {
		if _, present := eachEp[epoint]; present {
			return true
		}
	}
	return false
}
