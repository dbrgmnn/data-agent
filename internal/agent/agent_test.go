package agent_test

import (
	"data_agent/internal/agent"
	"flag"
	"os"
	"testing"
	"time"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantURL string
		wantDur time.Duration
		wantErr bool
	}{
		{
			name:    "valid flags",
			args:    []string{"cmd", "--url=amqp://guest:guest@localhost:5672/", "--interval=5"},
			wantURL: "amqp://guest:guest@localhost:5672/",
			wantDur: 5 * time.Second,
			wantErr: false,
		},
		{
			name:    "invalid URL scheme",
			args:    []string{"cmd", "--url=http://localhost:5672/"},
			wantErr: true,
		},
		{
			name:    "missing URL",
			args:    []string{"cmd"},
			wantErr: true,
		},
		{
			name:    "invalid interval defaults to 2",
			args:    []string{"cmd", "--url=amqp://guest:guest@localhost:5672/", "--interval=0"},
			wantURL: "amqp://guest:guest@localhost:5672/",
			wantDur: 2 * time.Second,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset and set flags for each test
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ExitOnError)
			os.Args = tt.args

			url, interval, err := agent.ParseFlags()

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr {
				if url != tt.wantURL {
					t.Errorf("expected URL %s, got %s", tt.wantURL, url)
				}
				if interval != tt.wantDur {
					t.Errorf("expected interval %v, got %v", tt.wantDur, interval)
				}
			}
		})
	}
}
