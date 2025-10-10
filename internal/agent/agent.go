package agent

import (
	"context"
	"flag"
	"log"
	"monitoring/internal/models"
	"regexp"
	"time"
)

// run the agent to collect and send metrics to RabbitMQ
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
			metricMsg := models.MetricMessage{
				Host:   CollectHostInfo(),
				Metric: CollectMetricInfo(),
			}
			if err := SendMetrics(&metricMsg, rabbitURL); err != nil {
				log.Println("Failed to send metrics:", err)
			}
			time.Sleep(2 * time.Second)
		}
	}
}

// parse flag --url
func ParseFlags() string {
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
