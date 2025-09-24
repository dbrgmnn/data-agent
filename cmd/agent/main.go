package main

import (
	"context"
	"flag"
	"log"
	"monitoring/internal/agent"
	"os"
	"os/signal"
	"regexp"
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
	url := parseFlags()
	agent.Run(ctx, url)
}

func parseFlags() string {
	// parse flag
	rabbitURL := flag.String("url", "", "RabbitMQ URL")
	flag.Parse()
	if *rabbitURL == "" {
		log.Fatal("RabbitMQ URL must be specified with --url")
	}
	// regex to validate URL format
	re := regexp.MustCompile(`^amqp://[^:]+:[^@]+@[^:]+:\d+/`)
	if !re.MatchString(*rabbitURL) {
		log.Fatal("RabbitMQ URL must match amqp://user:pass@host:port/")
	}
	return *rabbitURL
}
