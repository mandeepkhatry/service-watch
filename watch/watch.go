package watch

import (
	"encoding/json"
	"service-watch/internal/def"
	"service-watch/internal/models"

	"github.com/BurntSushi/toml"
)

//TODO Swagger Parsing, traversal... after Anish finishes deploying swagger

type ServiceWatcher struct {
	ConfigPath string
	AppConfig  models.AppConfig
}

//Init initializes ServiceWatcher.
func (s *ServiceWatcher) Init(configPath string) {
	s.ConfigPath = configPath
}

//ReadConfig decodes config file and initialize app config.
func (s *ServiceWatcher) ReadConfig() error {

	var config map[string]interface{}

	_, err := toml.DecodeFile(s.ConfigPath, &config)
	if err != nil {
		return err
	}

	bytesConfig, err := json.Marshal(config)
	if err != nil {
		return err
	}

	appConfig := models.AppConfig{}

	err = json.Unmarshal(bytesConfig, &appConfig)
	if err != nil {
		return err
	}

	s.AppConfig = appConfig

	return nil
}

//TODO Watch

//Watch watches overall working of the apis.
func (s *ServiceWatcher) Watch() error {

	if s.AppConfig.Api == nil {
		return def.ErrAppConfigUnregistered
	}

	return nil
}
