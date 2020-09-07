package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"service-watch/internal/def"
	"service-watch/internal/models"
	"strconv"
	"time"
)

type HTTPClient struct {
	Client *http.Client
	Host   string
	Port   int
}

//NewHTTPClient returns new HttpClient instance.
func NewHTTPClient(config map[string]interface{}) *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: time.Duration(config["timeout"].(int)) * time.Second,
		},
		Host: config["host"].(string),
		Port: config["timeout"].(int),
	}
}

func (c *HTTPClient) ExecuteRequest(method string, endpoint string, requestBody map[string]interface{}, requestConfig map[string]string) (models.HeartBeatResponse, error) {

	var req *http.Request
	var err error

	url := c.Host + ":" + strconv.Itoa(c.Port) + endpoint

	start := time.Now()

	if requestBody != nil {

		requestBytes, err := json.Marshal(requestBody)
		if err != nil {
			return models.HeartBeatResponse{}, err
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBytes))
		if err != nil {
			return models.HeartBeatResponse{}, err
		}

	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return models.HeartBeatResponse{}, err
		}
	}

	//Add Headers
	for eachConfigKey, eachConfigValue := range requestConfig {
		if _, headerConfirmed := def.RequestHeaders[eachConfigKey]; headerConfirmed {
			req.Header.Add(eachConfigKey, eachConfigValue)
		}
	}

	res, err := c.Client.Do(req)

	if err != nil {
		return models.HeartBeatResponse{}, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return models.HeartBeatResponse{}, err
	}

	var message interface{}
	err = json.Unmarshal(b, &message)
	if err != nil {
		return models.HeartBeatResponse{}, err
	}

	elapsedTime := time.Since(start).String()

	response := models.HeartBeatResponse{Status: res.Status, StatusCode: res.StatusCode, Message: message, ElapsedTime: elapsedTime}

	res.Body.Close()

	return response, nil

}
