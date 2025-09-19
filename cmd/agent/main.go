package main

import (
	"log"
	"monitoring/internal/agent"
	"time"
)

func main() {
	for {
		metric := agent.CollectMetrics()
		err := agent.SendMetrics(metric)
		if err != nil {
			log.Println("Failed to send metrics:", err)
		}
		time.Sleep(1 * time.Second)
	}
}
