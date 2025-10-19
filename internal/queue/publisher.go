package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"monitoring/internal/models"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

// create one connection and queue
func NewPublisher(server string) (*Publisher, error) {
	// connect to RabbitMQ server
	conn, err := amqp.Dial(server)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

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
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	log.Println("Connected to RabbitMQ and declare queue:", q.Name)
	return &Publisher{conn: conn, ch: ch, q: q}, nil
}

// publish a message without opening a new connection
func (p *Publisher) Publish(metricMsg *models.MetricMessage) error {
	// serialize metric to JSON
	body, err := json.Marshal(metricMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	// publish the message
	err = p.ch.Publish(
		"",       // exchange
		p.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
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

// close connection
func (p *Publisher) Close() {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	log.Println("Connection closed")
}
