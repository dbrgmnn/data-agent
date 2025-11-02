package agent_test

import (
	"data_agent/internal/agent"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestParseFlags_ValidInput checks parsing with valid URL and interval
func TestParseFlags_ValidInput(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--url", "amqp://guest:guest@localhost:5672/", "--interval", "5"}

	url, interval, err := agent.ParseFlags()

	assert.NoError(t, err)
	assert.Equal(t, "amqp://guest:guest@localhost:5672/", url)
	assert.Equal(t, 5*time.Second, interval)
}

// TestParseFlags_MissingURL verifies error when URL flag is missing
func TestParseFlags_MissingURL(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--interval", "2"}

	_, _, err := agent.ParseFlags()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be specified")
}

// TestParseFlags_InvalidURL checks error on invalid RabbitMQ URL scheme
func TestParseFlags_InvalidURL(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--url", "http://example.com", "--interval", "3"}

	_, _, err := agent.ParseFlags()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid RabbitMQ URL")
}

// TestParseFlags_DefaultInterval ensures default interval is used if 0 is provided
func TestParseFlags_DefaultInterval(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--url", "amqp://guest:guest@localhost:5672/", "--interval", "0"}

	_, interval, err := agent.ParseFlags()

	assert.NoError(t, err)
	assert.Equal(t, 2*time.Second, interval)
}
