package models

type HeartBeatResponse struct {
	Status      string      `json:"status"`
	StatusCode  int         `json:"code"`
	Message     interface{} `json:"message"`
	ElapsedTime string      `json:"elapsedtime"`
}
