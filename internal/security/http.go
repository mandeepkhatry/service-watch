package security

import (
	"service-watch/internal/client"
	"service-watch/internal/def"
	"service-watch/internal/utils"
)

type HTTPAuthenticationScheme struct {
	HttpClient  *client.HTTPClient
	Credentials map[string]interface{}
	Response    map[string]interface{}
}

func (s *HTTPAuthenticationScheme) Run() {

	res, _ := s.HttpClient.ExecuteRequest("POST", def.HTTPSecurityEndpoints[0], utils.ConvertMap(def.HTTPSecurityCredentials), nil)
	s.Response = res.Message.(map[string]interface{})
}
