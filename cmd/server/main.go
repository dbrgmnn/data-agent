package main

import (
	"log"
	"monitoring/internal/config"
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
	rabbitURL := config.LoadConfig().RabbitURLServer
	if err := server.StartMetricsConsumer(db, rabbitURL); err != nil {
		log.Fatal(err)
	}

	// block main
	select {}
}
