package models

import "github.com/getkin/kin-openapi/openapi3"

type AppConfig struct {
	Server     string
	Endpoints  map[string][]map[string]Endpoint
	Components openapi3.Components
}
