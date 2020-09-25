package models

type RearMostEndpoint struct {
	MethodName       string
	SpecificEndpoint string
	RequestBody      map[string]interface{}
	RequestConfig    map[string]string
}
