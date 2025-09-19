package server

import (
	"encoding/json"
	"log"
	"monitoring/internal/models"
	"net/http"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	var m models.Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Printf("Received metrics: %+v\n", m)

	w.WriteHeader(http.StatusOK)
}
