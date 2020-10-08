package utils

import (
	"service-watch/internal/def"
	"service-watch/internal/models"
	"time"
)

func LogExtract(ep string, method string, response models.HeartBeatResponse, responseValid bool) (string, map[string]interface{}) {

	status := ""

	if response.StatusCode == 0 {
		status = def.StatusTimeout
	} else {
		if !responseValid {
			status = def.StatusInvalidResponse
		} else {
			if _, present := def.AccpetedStatus[method][response.StatusCode]; present {
				status = def.StatusSuccess
			} else {
				status = def.StatusFailure
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
