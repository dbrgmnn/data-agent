package queue

import (
	"context"
	"data-agent/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	q      amqp.Queue
	ctx    context.Context
	server string
	mu     sync.Mutex
}

// create a new publisher with context
func NewPublisher(ctx context.Context, server string) *Publisher {
	return &Publisher{
		ctx:    ctx,
		server: server,
	}
}

// connect to RabbitMQ with mutex
func (p *Publisher) connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// return nil if connection closed
	if p.conn != nil && !p.conn.IsClosed() && p.ch != nil {
		return nil
	}

	// open connection
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
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn == nil || p.conn.IsClosed() {
		return fmt.Errorf("cannot publish, connection is closed or nil")
	}

	if p.ch == nil {
		return fmt.Errorf("cannot publish, channel is closed or nil")
	}

	// marshaling metrics
	body, err := json.Marshal(metricMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	// publish metrics
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
func (p *Publisher) StartMetricsPublisher() {
	for {
		if err := p.connect(); err != nil {
			log.Println("Publisher connection failed, retrying:", err)
			select {
			case <-p.ctx.Done():
				log.Println("Publisher stopped by context")
				return
			case <-time.After(5 * time.Second):
				continue
			}
		}

		p.mu.Lock()
		notifyClose := make(chan *amqp.Error, 1)
		p.conn.NotifyClose(notifyClose)
		p.mu.Unlock()

		select {
		case <-p.ctx.Done():
			log.Println("Publisher stopping...")
			p.Close()
			return
		case err := <-notifyClose:
			if err != nil {
				log.Println("Publisher connection closed, reconnecting:", err)
			}
		}
	}
}

// close channel and connection gracefully
func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			log.Println("Error closing channel:", err)
		}
		p.ch = nil
	}
	if p.conn != nil && !p.conn.IsClosed() {
		if err := p.conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}
	p.conn = nil
	log.Println("Publisher connection closed")
}
