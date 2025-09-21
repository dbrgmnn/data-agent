package main

import (
	"log"
	"monitoring/internal/server"
)

func main() {
	// Initialize database
	db, err := server.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Start RabbitMQ consumer
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	if err := server.StartMetricsConsumer(db, rabbitURL); err != nil {
		log.Fatal(err)
	}
	log.Println("Server is running and consuming metrics from RabbitMQ")

	// Block main
	select {}
}
