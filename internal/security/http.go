package security

import (
	"service-watch/internal/client"
	"service-watch/internal/utils"
)

type HTTPAuthenticationScheme struct {
	HttpClient        *client.HTTPClient
	Credentials       map[string]interface{}
	Response          map[string]interface{}
	SecurityEndpoints []string
}

func (s *HTTPAuthenticationScheme) Run() {

	res, _ := s.HttpClient.ExecuteRequest("POST", s.SecurityEndpoints[0], utils.ConvertMap(s.Credentials), nil)
	s.Response = res.Message.(map[string]interface{})
}
