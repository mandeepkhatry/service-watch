package parser

import (
	"fmt"
	"service-watch/internal/models"
	"service-watch/internal/utils"
	"strings"
)

func ParseEndpoint(endpoint string) []string {

	epParts := strings.Split(endpoint, "/")

	roots := make([]string, 0)

	if epParts[0] == "" {
		epParts = epParts[1:]
	}

	for i := range epParts {

		roots = append(roots, "/"+strings.Join(epParts[0:i+1], "/"))
	}

	return roots
}

func GenerateSpecificEndpoint(ep string, endpointsDataBuffer map[string]map[string]models.DataBuffer) string {

	epParts := strings.Split(ep, "/")

	finalEp := "/"

	oneTimePost := true

	current := 0

	epParts = epParts[1:]

	for i, epPart := range epParts {
		if utils.ValidateResource(epPart) {

			root := finalEp + strings.Join(epParts[current:i], "/")

			resource := ""

			if oneTimePost {
				//epPart is resource since already validated
				resource = endpointsDataBuffer[root]["POST"].Response["data"].(map[string]interface{})[utils.SeperateResource(epPart)].(string)
				oneTimePost = false
			} else {
				resource = endpointsDataBuffer[root]["GET"].Response["data"].(map[string]interface{})[utils.SeperateResource(epPart)].(string)
			}

			finalEp = root + "/" + resource + "/"

			if i <= len(epPart)-1 {
				current = i + 1
			}

		}

	}

	if finalEp == "/" {
		finalEp = ep
	}

	if !utils.ValidateResource(epParts[len(epParts)-1]) {
		finalEp += epParts[len(epParts)-1]
	} else {
		finalEp = strings.TrimSuffix(finalEp, "/")
	}

	fmt.Println("FINAL :", finalEp)

	return finalEp

}
