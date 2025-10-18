package main

import (
	"context"
	"log"
	"monitoring/internal/config"
	initDB "monitoring/internal/db"
	q "monitoring/internal/queue"
	"os"
	"os/signal"
	"syscall"
)

// main function to start the RabbitMQ consumer
func main() {
	// initialize database
	db, err := initDB.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// load configuration
	cfg := config.LoadConfig()
	rabbitURL := cfg.RabbitURL

	// create a context that is canceled on exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle termination signals in a separate goroutine
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Shutting down ingestor...")
		cancel()
	}()

	// start RabbitMQ consumer
	log.Println("Starting RabbitMQ consumer...")
	if err := q.StartMetricsConsumer(ctx, db, rabbitURL); err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
	}
}
