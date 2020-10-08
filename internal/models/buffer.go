package models

type DataBuffer struct {
	Request  map[string]interface{}
	Response map[string]interface{}
}

func (d *DataBuffer) AssignRequest(requestData map[string]interface{}) {
	d.Request = requestData
}

func (d *DataBuffer) AssignResponse(responseData map[string]interface{}) {
	d.Response = responseData
}
