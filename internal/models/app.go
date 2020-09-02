package models

import "github.com/xeipuuv/gojsonschema"

type AppConfig struct {
	Api map[string]*ApiConfig `json:"api"`
}

type ApiConfig struct {
	Endpoint string                   `json:"endpoint"`
	Methods  map[string]*MethodConfig `json:"methods"`
}

type RequestConfig struct {
	Auth        string                 `json:"auth"`
	ContentType string                 `json:"content_type"`
	Validator   map[string]interface{} `json:"validator"`
}

type ResponseConfig struct {
	ContentType string `json:"content_type"`
}

type MethodConfig struct {
	Request  *RequestConfig  `json:"request"`
	Response *ResponseConfig `json:"response"`
	Schema   *gojsonschema.Schema
}
