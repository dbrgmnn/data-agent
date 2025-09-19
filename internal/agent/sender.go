package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"monitoring/internal/models"
	"net/http"
)

func SendMetrics(metric models.Metric) error {
	data, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	resp, err := http.Post("http://localhost:8080/metrics", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	return nil
}
