package main

import (
	"flag"
	"log"
	"monitoring/internal/agent"
)

func main() {
	serverFlag := flag.String("server", "", "Server URL to send metrics to")
	flag.Parse()

	serverURl := *serverFlag
	if serverURl == "" {
		serverURl = agent.DefaultServerURL()
	}

	log.Println("Sending metrics to:", serverURl)

	agent.Run(serverURl)
}
