package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	save "monitoring/internal/db"
	"monitoring/internal/models"

	"github.com/streadway/amqp"
)

// saves a metric to the database
func StartMetricsConsumer(ctx context.Context, db *sql.DB, rabbitURL string) error {
	// connect to RabbitMQ server
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// declare a queue
	q, err := ch.QueueDeclare(
		"metrics", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// subscribe to the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}

	// process messages
	log.Println("Started metrics consumer")

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down metrics consumer")
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
				d.Nack(false, false) // don't send to queue
				continue
			}

			// send metric to database
			if err := save.SaveMetric(db, &metric); err != nil {
				log.Println("Failed to save metric:", err)
				d.Nack(false, true) // send to queue again
				continue
			}

			d.Ack(false) // acknowledge message
			log.Printf("Metric saved from queue: %+v\n", metric)
		}
	}
}
