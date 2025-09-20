package server

import (
	"encoding/json"
	"log"
	"monitoring/internal/models"
	"net/http"
)

func MetricsHandler(writer http.ResponseWriter, requesr *http.Request) {
	var metric models.Metric
	err := json.NewDecoder(requesr.Body).Decode(&metric)
	if err != nil {
		http.Error(writer, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Printf("Received metrics: %+v\n", metric)

	writer.WriteHeader(http.StatusOK)
}
