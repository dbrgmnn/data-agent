package agent

import (
	"context"
	"flag"
	"log"
	"monitoring/internal/models"
	rabbit "monitoring/internal/queue"
	"net/url"
	"time"
)

// run the agent to collect and send metrics to RabbitMQ
func Run(ctx context.Context, rabbitURL string, interval time.Duration) {
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
			host, err := CollectHostInfo()
			if err != nil {
				log.Println("CollectHostInfo error:", err)
				continue
			}
			metric, err := CollectMetricInfo()
			if err != nil {
				log.Println("CollectMetricInfo error:", err)
				continue
			}
			metricMsg := models.NewMetricMessage(&host, &metric)

			log.Printf("Sending metrics from [%s] to [%s]:", host.Hostname, rabbitURL)
			if err := rabbit.SendMetrics(metricMsg, rabbitURL); err != nil {
				log.Printf("Failed to send metrics: %v\n", err)
			}
		}
	}
}

// parse flag --url and --interval
func ParseFlags() (string, time.Duration) {
	rabbitURL := flag.String("url", "", "RabbitMQ URL")
	interval := flag.Int("interval", 2, "Interval in seconds between metric collections")
	flag.Parse()

	if *interval <= 0 {
		log.Println("Invalid interval, using default 2 seconds")
		*interval = 2
	}

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
