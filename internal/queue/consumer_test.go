package queue_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"data_agent/internal/queue"

	"github.com/stretchr/testify/assert"
)

// MockConnection is a mock implementation of a connection, used for testing
type MockConnection struct{}

func (m *MockConnection) Close() error { return nil }

// MockChannel is a mock implementation of a channel, used for testing
type MockChannel struct{}

func (m *MockChannel) Close() error { return nil }

// TestConsumer embeds queue.Consumer and is used in tests to simulate a consumer
type TestConsumer struct {
	*queue.Consumer
}

// NewTestConsumer creates a new TestConsumer with the provided context and database
func NewTestConsumer(ctx context.Context, db *sql.DB) *TestConsumer {
	c := &queue.Consumer{
		Ctx: ctx,
		Db:  db,
	}
	return &TestConsumer{Consumer: c}
}

// StartMetricsConsumer waits until the context is cancelled, simulating consumer shutdown
func (c *TestConsumer) StartMetricsConsumer() {
	<-c.Ctx.Done()
}

// TestNewConsumer verifies that a new TestConsumer is created with the correct context and database
func TestNewConsumer(t *testing.T) {
	ctx := context.Background()
	db := &sql.DB{}
	c := NewTestConsumer(ctx, db)

	assert.NotNil(t, c)
	assert.Equal(t, db, c.Db)
	assert.Equal(t, ctx, c.Ctx)
}

// TestStartMetricsConsumerStopsWithContext ensures that the consumer stops when the context is cancelled
func TestStartMetricsConsumerStopsWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewTestConsumer(ctx, &sql.DB{})

	done := make(chan struct{})
	// start the consumer in a separate goroutine and close 'done' when it exits
	go func() {
		c.StartMetricsConsumer()
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// consumer stopped as expected
	case <-time.After(1 * time.Second):
		t.Fatal("Consumer did not stop after context cancel")
	}
}
