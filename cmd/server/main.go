package main

import (
	"log"
	"monitoring/internal/server"
)

func main() {
	// initialize database
	db, err := server.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// start RabbitMQ consumer
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	if err := server.StartMetricsConsumer(db, rabbitURL); err != nil {
		log.Fatal(err)
	}
	log.Println("Server is running and consuming metrics from RabbitMQ")

	// block main
	select {}
}
