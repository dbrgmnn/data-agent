package main

import (
	"context"
	"data-agent/internal/config"
	dataBase "data-agent/internal/db"
	"data-agent/internal/queue"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// main function to start the RabbitMQ consumer
func main() {
	// add prefix for logs
	log.SetPrefix("[ingestor] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initialize database
	db, err := dataBase.InitDB()
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
		log.Println("Stopping ingestor...")
		cancel()
	}()

	// create and start consumer
	consumer := queue.NewConsumer(ctx, db, rabbitURL)
	consumer.StartMetricsConsumer()
}
