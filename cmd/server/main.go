package main

import (
	"log"
	"monitoring/internal/server"
	"net/http"
)

func main() {
	db, err := server.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/metrics", server.MetricsHandler(db))
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
