package main

import (
	"context"
	"data-agent/internal/agent"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// main function to run the agent
func main() {
	// add prefix for logs
	log.SetPrefix("[agent] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// create a context that is canceled on exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle termination signals in a separate goroutine
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Stopping agent...")
		cancel()
	}()

	// parse flags and run the agent
	url, interval := agent.ParseFlags()
	agent.Run(ctx, url, interval)
}
