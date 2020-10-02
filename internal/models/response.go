package models

type HeartBeatResponse struct {
	Status                string                 `json:"status"`
	StatusCode            int                    `json:"code"`
	Message               interface{}            `json:"message"`
	ElapsedTime           string                 `json:"elapsedtime"`
	Timeout               bool                   `json:"timeout"`
	ResponseValidity      bool                   `json:"responsevalidity"`
	ResponseValidityError map[string]interface{} `json:"responsevalidityerror"`
}
