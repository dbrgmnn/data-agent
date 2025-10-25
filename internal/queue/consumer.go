package queue

import (
	"context"
	dataBase "data_agent/internal/db"
	"data_agent/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	db        *sql.DB
	ctx       context.Context
	rabbitURL string
}

// create a new consumer with context
func NewConsumer(ctx context.Context, db *sql.DB, rabbitURL string) *Consumer {
	return &Consumer{
		db:        db,
		ctx:       ctx,
		rabbitURL: rabbitURL,
	}
}

// connect to RabbitMQ
func (c *Consumer) connect() error {
	// open connection
	conn, err := amqp.DialConfig(c.rabbitURL, amqp.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
	})
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// open channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	c.conn = conn
	c.ch = ch
	return nil
}

// saves a metric to the database
func (c *Consumer) consumeMetrics() error {
	// declare a queue
	q, err := c.ch.QueueDeclare("metrics", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// subscribe to the queue
	msgs, err := c.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to the queue: %w", err)
	}

	// process messages
	for {
		select {
		case <-c.ctx.Done():
			return nil

		case d, ok := <-msgs:
			if !ok {
				log.Println("Message channel closed")
				return nil
			}

			// decode message
			var metric models.MetricMessage
			if err := json.Unmarshal(d.Body, &metric); err != nil {
				log.Println("Failed to decode metric:", err)
				// don't send to queue
				d.Nack(false, false)
				continue
			}

			// send metric to database
			if err := dataBase.SaveMetric(c.ctx, c.db, &metric); err != nil {
				log.Println("Failed to save metric:", err)
				// send to queue again
				d.Nack(false, true)
				continue
			}

			// acknowledge message
			d.Ack(false)
			log.Printf("Metric saved from queue: host=%s", metric.Host.Hostname)
		}
	}
}

// consume metrics
func (c *Consumer) StartMetricsConsumer() {
	for {
		if err := c.connect(); err != nil {
			log.Println("Consumer connection failed, retrying in 5s:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Connected to RabbitMQ")
		if err := c.consumeMetrics(); err != nil {
			log.Println("Consume error, reconnecting:", err)
		}

		select {
		case <-c.ctx.Done():
			log.Println("Consumer stopped by context")
			c.Close()
			return
		case <-time.After(5 * time.Second):
		}
	}
}

// close channel and connection gracefully
func (c *Consumer) Close() {
	if c.ch != nil {
		if err := c.ch.Close(); err != nil {
			log.Println("Error closing channel:", err)
		}
		c.ch = nil
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
		c.conn = nil
	}
	log.Println("Consumer connection closed")
}
