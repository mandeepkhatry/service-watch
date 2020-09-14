package models

type AppConfig struct {
	Server    string
	Endpoints map[string][]map[string]Endpoint
}
