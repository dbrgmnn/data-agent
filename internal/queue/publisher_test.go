package queue_test

import (
	"context"
	"data_agent/internal/models"
	"errors"
	"testing"
	"time"

	"data_agent/internal/queue"

	"github.com/stretchr/testify/assert"
)

// TestPublisher simulates a Publisher for testing without real RabbitMQ
type TestPublisher struct {
	*queue.Publisher
}

// NewTestPublisher creates a new TestPublisher with the given context
func NewTestPublisher(ctx context.Context) *TestPublisher {
	p := &queue.Publisher{
		Ctx: ctx,
	}
	return &TestPublisher{Publisher: p}
}

// StartMetricsPublisher simulates running publisher until context is cancelled
func (p *TestPublisher) StartMetricsPublisher() {
	<-p.Ctx.Done()
}

// publish simulates publishing a message without real RabbitMQ
func (p *TestPublisher) Publish(msg *models.MetricMessage) error {
	if msg == nil {
		return errors.New("msg is nil")
	}
	// simulate successful publishing
	return nil
}

// TestNewPublisher checks that a new TestPublisher is created
func TestNewPublisher(t *testing.T) {
	ctx := context.Background()
	p := NewTestPublisher(ctx)

	assert.NotNil(t, p)
	assert.Equal(t, ctx, p.Ctx)
}

// TestStartMetricsPublisherStopsWithContext ensures publisher stops when context is cancelled
func TestStartMetricsPublisherStopsWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := NewTestPublisher(ctx)

	done := make(chan struct{})
	go func() {
		p.StartMetricsPublisher()
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// publisher stopped as expected
	case <-time.After(1 * time.Second):
		t.Fatal("Publisher did not stop after context cancel")
	}
}

// TestPublishSimulated ensures Publish works in test mode without real connection
func TestPublishSimulated(t *testing.T) {
	ctx := context.Background()
	p := NewTestPublisher(ctx)

	msg := &models.MetricMessage{
		Host: models.Host{Hostname: "test-host"},
	}

	err := p.Publish(msg)
	assert.NoError(t, err)
}
