package utils

import (
	"service-watch/internal/def"
	"service-watch/internal/models"
)

func LogExtract(ep string, method string, response models.HeartBeatResponse) (string, map[string]interface{}) {

	status := ""

	if response.StatusCode == 0 {
		status = "timeout"
	} else {
		if _, present := def.AccpetedStatus[method][response.StatusCode]; present {
			status = "success"
		} else {
			status = "failure"
		}
	}

	return status, map[string]interface{}{
		"endpoint":    ep,
		"method":      method,
		"status":      response.Status,
		"code":        response.StatusCode,
		"message":     response.Message,
		"elapsedtime": response.ElapsedTime,
		"timeout":     response.Timeout,
	}
}
