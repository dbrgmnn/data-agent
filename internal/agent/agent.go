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
	// create one connection
	publisher, err := rabbit.NewPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}
	defer publisher.Close()

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
			if err := collectAndSend(publisher); err != nil {
				log.Println("Failed to send metrics:", err)
			}
		}
	}
}

// collect and send metrics
func collectAndSend(publisher *rabbit.Publisher) error {
	host, err := CollectHostInfo()
	if err != nil {
		return err
	}
	metric, err := CollectMetricInfo()
	if err != nil {
		return err
	}
	metricMsg := models.NewMetricMessage(&host, &metric)
	return publisher.Publish(metricMsg)
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
