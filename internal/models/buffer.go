package models

type DataBuffer struct {
	Request  map[string]interface{}
	Response map[string]interface{}
}

func (d *DataBuffer) AssignRequest(requestData map[string]interface{}) {
	d.Request = make(map[string]interface{})
	for k, v := range requestData {
		d.Request[k] = v
	}
}

func (d *DataBuffer) AssignResponse(responseData map[string]interface{}) {
	d.Response = make(map[string]interface{})
	for k, v := range responseData {
		d.Response[k] = v
	}
}