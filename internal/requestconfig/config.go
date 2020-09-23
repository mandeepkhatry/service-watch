package requestconfig

type RequestConfig struct {
	Content map[string]string
}

func NewRequestConfig(config map[string]interface{}) *RequestConfig {

	content := make(map[string]string)

	if _, present := config["data"]; present {
		content["Authorization"] = "Bearer " + config["data"].(map[string]interface{})["access_token"].(string)
	}

	return &RequestConfig{
		Content: content,
	}
}
