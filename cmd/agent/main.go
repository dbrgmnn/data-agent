package main

import (
	"context"
	"log"
	"monitoring/internal/agent"
	"monitoring/internal/config"
	"os"
	"os/signal"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Context
	ctx, cancel := context.WithCancel(context.Background())

	// Handel signal to stop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Ctrl+C detected, stopping...")
		cancel() // cancel context
	}()
	// Run agent
	agent.Run(ctx, cfg.RabbitURL)
}
