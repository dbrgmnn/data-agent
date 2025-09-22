package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"monitoring/internal/models"

	"github.com/streadway/amqp"
)

func StartMetricsConsumer(db *sql.DB, rabbitURL string) error {
	// connect to RabbitMQ server
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

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

	// start a goroutine to process messages
	go func() {
		for d := range msgs {
			var metric models.ExtendedMetrics
			// decode message
			if err := json.Unmarshal(d.Body, &metric); err != nil {
				log.Println("Failed to decode metric:", err)
				d.Nack(false, false) // don't send to queue
				continue
			}

			// send metric to database
			if err := SaveMetric(db, metric); err != nil {
				log.Println("Failed to save metric:", err)
				d.Nack(false, true) // send to queue again
				continue
			}

			d.Ack(false) // acknowledge message
			log.Printf("Metric seved from queue: %+v\n", metric)
		}
	}()

	log.Println("Started metrics consumer")
	return nil
}
