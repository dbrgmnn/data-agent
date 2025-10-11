package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"monitoring/internal/models"

	"github.com/streadway/amqp"
)

// send metrics to RabbitMQ
func SendMetrics(metricMsg *models.MetricMessage, server string) error {
	// connect to RabbitMQ server
	conn, err := amqp.Dial(server)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// declare a queue
	q, err := ch.QueueDeclare(
		"metrics", // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// serialize metric to JSON
	body, err := json.Marshal(metricMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	// publish the message
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // make message persistent
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Metric sent to queue: host=%s time=%s", metricMsg.Host.Hostname, metricMsg.Metric.Time)
	return nil
}
