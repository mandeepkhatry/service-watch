package models

import "github.com/getkin/kin-openapi/openapi3"

type Endpoint struct {
	Root    string
	Methods []map[string]*openapi3.Operation
}
