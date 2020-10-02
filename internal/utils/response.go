package utils

import (
	"service-watch/internal/def"
	"service-watch/internal/models"
	"time"
)

func LogExtract(ep string, method string, response models.HeartBeatResponse, responseValid bool) (string, map[string]interface{}) {

	status := ""

	if response.StatusCode == 0 {
		status = "timeout"
	} else {
		if !responseValid {
			status = "invalid_response"
		} else {
			if _, present := def.AccpetedStatus[method][response.StatusCode]; present {
				status = "success"
			} else {
				status = "failure"
			}
		}

	}

	timestamp := time.Now()

	return status, map[string]interface{}{
		"endpoint":    ep,
		"method":      method,
		"status":      response.Status,
		"code":        response.StatusCode,
		"message":     response.Message,
		"elapsedtime": response.ElapsedTime,
		"timeout":     response.Timeout,
		"timestamp":   timestamp,
	}
}
