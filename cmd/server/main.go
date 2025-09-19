package main

import (
	"log"
	"monitoring/internal/server"
	"net/http"
)

func main() {
	http.HandleFunc("/metrics", server.MetricsHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
