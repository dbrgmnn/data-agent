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
	// load configuration
	cfg := config.LoadConfig()

	// context
	ctx, cancel := context.WithCancel(context.Background())

	// handle signal to stop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Ctrl+C detected, stopping...")
		cancel() // cancel context
	}()
	// run agent
	agent.Run(ctx, cfg.RabbitURL)
}
