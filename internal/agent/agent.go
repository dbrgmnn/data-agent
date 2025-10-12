package agent

import (
	"context"
	"flag"
	"log"
	"monitoring/internal/models"
	q "monitoring/internal/queue"
	"net/url"
	"time"
)

// run the agent to collect and send metrics to RabbitMQ
func Run(ctx context.Context, rabbitURL string, interval time.Duration) {
	log.Println("Sending metrics to RabbitMQ:", rabbitURL)

	// send metrics every N seconds
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		// check if context is done
		select {
		case <-ctx.Done():
			log.Println("Agent stopped")
			return
		case <-ticker.C:
			// collect and send metrics
			host := CollectHostInfo()
			metric := CollectMetricInfo()
			metricMsg := models.NewMetricMessage(&host, &metric)

			if err := q.SendMetrics(metricMsg, rabbitURL); err != nil {
				log.Println("Failed to send metrics:", err)
			}
		}
	}
}

// parse flag --url and --interval
func ParseFlags() (string, time.Duration) {
	rabbitURL := flag.String("url", "", "RabbitMQ URL")
	interval := flag.Int("interval", 2, "Interval in seconds between metric collections")
	flag.Parse()

	if *rabbitURL == "" {
		log.Fatal("RabbitMQ URL must be specified with --url")
	}
	// validate URL format
	u, err := url.Parse(*rabbitURL)
	if err != nil || u.Scheme != "amqp" {
		log.Fatalf("Invalid RabbitMQ URL: %s (expected format amqp://user:pass@host:port/)", *rabbitURL)
	}
	return *rabbitURL, time.Duration(*interval) * time.Second
}
