package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// LogLevel can be "info", "warn", "error"
func SendLog(service string, requestID string, level string, message string, extra map[string]interface{}) {
	payload := map[string]interface{}{
		"service":    service,
		"request_id": requestID,
		"level":      level,
		"message":    message,
	}

	// merge extra fields if provided
	for k, v := range extra {
		payload[k] = v
	}

	data, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 2 * time.Second}
	req, _ := http.NewRequest("POST", "http://log-service:9000/logs", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client.Do(req) // best effort, ignore errors
}
