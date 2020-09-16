package utils

import (
	"strings"
)

func FindComponent(ref string) (string, string) {
	references := strings.Split(ref, "/")

	if len(references) > 1 {
		return references[2], references[3]
	}
	return "", ""
}
