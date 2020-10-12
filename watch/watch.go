package watch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"service-watch/internal/def"
	"service-watch/internal/heartbeat"
	"service-watch/internal/loader"
	"service-watch/internal/logs"
	"service-watch/internal/models"
	"service-watch/internal/server"
	"service-watch/internal/store"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

type ServiceWatcher struct {
	ApiConfiguration    models.AppConfig
	Timeout             int
	LogsDir             string
	Store               string
	Periodicity         string
	LogsFlush           string
	SecurityEndpoints   []string
	SecurityCredentials map[string]interface{}
}

//NewServiceWatcher returns new ServiceWatcher instance.
func NewServiceWatcher(configPath string) (*ServiceWatcher, error) {
	configFile, err := os.Open(configPath)

	if err != nil {
		return &ServiceWatcher{}, err
	}

	var watchConfig map[string]interface{}

	watchConfigByte, _ := ioutil.ReadAll(configFile)

	json.Unmarshal(watchConfigByte, &watchConfig)

	configFile.Close()

	//Comment for now...

	// ep := watchConfig["host"].(string) + watchConfig["endpoint"].(string)

	// resp, err := http.Get(ep)

	// if err != nil {
	// 	return &ServiceWatcher{}, err
	// }

	// body, err := ioutil.ReadAll(resp.Body)

	// if err != nil {
	// 	return &ServiceWatcher{}, err
	// }

	// resp.Body.Close()

	openApiFile, err := os.Open("config/test.json")

	if err != nil {
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(openApiFile)

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(byteValue)
	if err != nil {
		return &ServiceWatcher{}, err
	}

	openApiFile.Close()

	appConfig, err := loader.LoadSwagger(swagger, watchConfig)

	if err != nil {
		return &ServiceWatcher{}, err
	}

	if _, storePresent := store.Stores[watchConfig["store"].(string)]; !storePresent {
		return &ServiceWatcher{}, def.ErrStoreUnavailable
	}

	securityEndpoints := make([]string, 0)

	if _, present := watchConfig["security_endpoints"]; present {
		for _, v := range watchConfig["security_endpoints"].([]interface{}) {
			securityEndpoints = append(securityEndpoints, v.(string))
		}
	}

	credentials := make(map[string]interface{})

	if _, present := watchConfig["credentials"]; present {
		for k, v := range watchConfig["credentials"].(map[string]interface{}) {
			credentials[k] = v
		}

	}

	return &ServiceWatcher{
		ApiConfiguration:    *appConfig,
		Timeout:             int(watchConfig["timeout"].(float64)),
		LogsDir:             watchConfig["logs_dir"].(string),
		Store:               watchConfig["store"].(string),
		Periodicity:         watchConfig["periodicity"].(string),
		LogsFlush:           watchConfig["logs_flush"].(string),
		SecurityEndpoints:   securityEndpoints,
		SecurityCredentials: credentials,
	}, nil

}

//ValidateAppSpecificRequirements validates app specific requirements.
func (s *ServiceWatcher) ValidateAppSpecificRequirements() error {

	if reflect.DeepEqual(s.ApiConfiguration, models.AppConfig{}) {
		return def.ErrSwaggerConfigUnregistered
	}

	if len(s.ApiConfiguration.Server) == 0 {
		return def.ErrServersUnregistered
	}

	return nil

}

//Watch watches overall working of the apis.
func (s *ServiceWatcher) Watch() error {

	err := s.ValidateAppSpecificRequirements()

	if err != nil {
		return err
	}

	serverURL := s.ApiConfiguration.Server
	var config = map[string]interface{}{
		"host":    serverURL,
		"timeout": s.Timeout,
	}

	periodicity, err := time.ParseDuration(s.Periodicity)

	if err != nil {
		return err
	}

	logsFlushPeriod, err := time.ParseDuration(s.LogsFlush)

	if err != nil {
		return err
	}

	ticker := time.NewTicker(periodicity)

	storeLog := logs.NewLog(s.Store, s.LogsDir)

	var logsFlushTimeStart time.Time
	logsFlushStatus := true
	go server.RunServer(storeLog)

	for range ticker.C {

		if logsFlushStatus {
			logsFlushTimeStart = time.Now()
			logsFlushStatus = false
		}

		fmt.Println("************************************************************")
		log, err := heartbeat.ProcessRequest(s.ApiConfiguration, config, s.SecurityEndpoints, s.SecurityCredentials)

		if err != nil {
			return err
		}

		storeLog.StoreLogs(log)
		fmt.Println("***********************************************************")

		if time.Since(logsFlushTimeStart) >= logsFlushPeriod {
			//flush logs store

			fmt.Println("[flush] Store Logs")
			err := storeLog.FlushLogs()
			if err != nil {
				return err
			}

			//flip status
			logsFlushStatus = true

		}

	}

	return nil
}
