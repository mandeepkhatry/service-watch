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

func IsSpecificItem(ep string) (string, string, bool) {
	epParts := strings.Split(ep, "/")

	lastPart := epParts[len(epParts)-1]

	isResource := ValidateResource(lastPart)

	if isResource {
		return strings.Join(epParts[0:len(epParts)-1], "/"), lastPart[1 : len(lastPart)-1], isResource
	}

	return "", "", isResource
}

func SeperateResource(epPart string) string {
	s := strings.TrimPrefix(epPart, "{")
	return strings.TrimSuffix(s, "}")
}
