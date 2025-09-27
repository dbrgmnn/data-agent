package main

import (
	"context"
	"log"
	"monitoring/internal/agent"
	"os"
	"os/signal"
)

func main() {
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
	url := agent.ParseFlags()
	agent.Run(ctx, url)
}
