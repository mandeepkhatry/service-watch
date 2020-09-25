package content

import "github.com/getkin/kin-openapi/openapi3"

func GetContent(content openapi3.Content) string {
	contentType := ""
	for c := range content {
		contentType = c
		break
	}
	return contentType
}
