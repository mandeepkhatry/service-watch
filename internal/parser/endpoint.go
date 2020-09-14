package parser

import (
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
