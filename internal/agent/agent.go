package agent

import (
	"context"
	"log"
	"time"
)

func Run(ctx context.Context, rabbitURL string) {
	log.Println("Sending metrics to RabbitMQ:", rabbitURL)
	for {
		// check if context is done
		select {
		case <-ctx.Done():
			log.Println("Stoping agent")
			return
		default:
			// collect and send metrics every 5 seconds
			metric := CollectMetrics()
			err := SendMetrics(metric, rabbitURL)
			if err != nil {
				log.Println("Failed to send metrics:", err)
			}
			time.Sleep(5 * time.Second)

		}
	}
}
