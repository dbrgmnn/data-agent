package main

import (
	"context"
	"log"
	"monitoring/internal/agent"
	"os"
	"os/signal"
	"syscall"
)

// main function to run the agent
func main() {
	// create a context that is canceled on exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle signal to stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Stopping agent...")
		cancel() // cancel context
	}()
	// run agent
	url, interval := agent.ParseFlags()
	agent.Run(ctx, url, interval)
}
