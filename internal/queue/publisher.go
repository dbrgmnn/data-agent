package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"monitoring/internal/models"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	q      amqp.Queue
	ctx    context.Context
	server string
}

// create a new publisher with context
func NewPublisher(ctx context.Context, server string) *Publisher {
	return &Publisher{
		ctx:    ctx,
		server: server,
	}
}

// connect to RabbitMQ
func (p *Publisher) connect() error {
	conn, err := amqp.DialConfig(p.server, amqp.Config{
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

	q, err := ch.QueueDeclare("metrics", true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	p.conn = conn
	p.ch = ch
	p.q = q
	log.Println("Publisher connected and queue declared:", q.Name)
	return nil
}

// publish a message without opening a new connection
func (p *Publisher) Publish(metricMsg *models.MetricMessage) error {
	// open new channel when closed
	if p.ch == nil {
		if err := p.connect(); err != nil {
			return err
		}
	}

	// marshaling metrics
	body, err := json.Marshal(metricMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	err = p.ch.Publish("", p.q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return fmt.Errorf("failed to publish metric: %w", err)
	}

	log.Printf("Metric sent to queue: host=%s", metricMsg.Host.Hostname)
	return nil
}

// publish metrics
func (p *Publisher) StartPublisher() {
	for {
		if err := p.connect(); err != nil {
			log.Println("Publisher connection failed, retrying in 5s:", err)
			select {
			case <-time.After(5 * time.Second):
				continue
			case <-p.ctx.Done():
				log.Println("Publisher stopped by context")
				return
			}
		}

		notifyClose := make(chan *amqp.Error)
		p.conn.NotifyClose(notifyClose)

		select {
		case <-p.ctx.Done():
			log.Println("Publisher stopping...")
			p.Close()
			return
		case err := <-notifyClose:
			log.Println("Publisher connection closed, reconnecting:", err)
			p.Close()
			time.Sleep(5 * time.Second)
		}
	}
}

// close connection
func (p *Publisher) Close() {
	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			log.Println("Error closing channel:", err)
		}
		p.ch = nil
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
		p.conn = nil
	}
	log.Println("Publisher connection closed")
}
