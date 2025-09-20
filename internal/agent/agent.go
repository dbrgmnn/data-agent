package agent

import (
	"log"
	"os"
	"time"
)

func DefaultServerURL() string {
	url := os.Getenv("SERVER_URL")
	if url == "" {
		url = "http://localhost:8080/metrics"
	}
	return url
}

func Run(serverURL string) {
	for {
		metric := CollectMetrics()
		err := SendMetrics(metric, serverURL)
		if err != nil {
			log.Println("Failed to send metrics:", err)
		}
		time.Sleep(5 * time.Second)
	}
}
