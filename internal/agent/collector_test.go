package agent_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"data_agent/internal/agent"
)

// TestCollectHostInfo verifies that CollectHostInfo returns valid host information without error
func TestCollectHostInfo(t *testing.T) {
	hostInfo, err := agent.CollectHostInfo()
	assert.NoError(t, err, "CollectHostInfo should not return an error")
	assert.NotNil(t, hostInfo, "HostInfo should not be nil")
	assert.NotEmpty(t, hostInfo.Hostname, "Hostname should not be empty")
}

// TestCollectMetricInfo verifies that CollectMetricInfo returns valid metric information without error
func TestCollectMetricInfo(t *testing.T) {
	metricInfo, err := agent.CollectMetricInfo()
	assert.NoError(t, err, "CollectMetricInfo should not return an error")
	assert.NotNil(t, metricInfo, "MetricInfo should not be nil")
}

// TestCollectDiskMetric verifies that CollectDiskMetric returns disk metrics without error
func TestCollectDiskMetric(t *testing.T) {
	diskMetric, err := agent.CollectDiskMetric()
	assert.NoError(t, err, "CollectDiskMetric should not return an error")
	assert.NotNil(t, diskMetric, "DiskMetric should not be nil")
}

// TestCollectNetMetric verifies that CollectNetMetric returns network metrics without error
func TestCollectNetMetric(t *testing.T) {
	netMetric, err := agent.CollectNetMetric()
	assert.NoError(t, err, "CollectNetMetric should not return an error")
	assert.NotNil(t, netMetric, "NetMetric should not be nil")
}
