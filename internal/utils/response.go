package utils

import "service-watch/internal/models"

func LogExtract(ep string, response models.HeartBeatResponse) map[string]interface{} {
	return map[string]interface{}{
		"endpoint":    ep,
		"status":      response.Status,
		"code":        response.StatusCode,
		"message":     response.Message,
		"elapsedtime": response.ElapsedTime,
		"timeout":     response.Timeout,
	}
}
