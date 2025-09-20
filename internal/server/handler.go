package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"monitoring/internal/models"
	"net/http"
)

func MetricsHandler(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, reader *http.Request) {
		var metric models.Metric
		if err := json.NewDecoder(reader.Body).Decode(&metric); err != nil {
			http.Error(writer, "Invalid JSON", http.StatusBadRequest)
			return
		}

		log.Printf("Received metrics: %+v\n", metric)

		if err := SaveMetric(db, metric); err != nil {
			http.Error(writer, "Failed to save metrics", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
