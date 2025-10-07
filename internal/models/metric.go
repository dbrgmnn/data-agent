package models

import "time"

// represents a monitored host system
type Host struct {
	ID          int64  `json:"id"`
	Hostname    string `json:"hostname"`
	OS          string `json:"os"`
	Platform    string `json:"platform"`
	PlatformVer string `json:"platformver"`
	KernelVer   string `json:"kernelver"`
}

// base metrics for all systems
type Metric struct {
	ID      int64        `json:"id"`
	HostID  int64        `json:"host_id"`
	Uptime  uint64       `json:"uptime"`
	CPU     float64      `json:"cpu"`
	RAM     float64      `json:"ram"`
	Disk    []DiskMetric `json:"disk,omitempty"`
	Network []NetMetric  `json:"network,omitempty"`
	Time    time.Time    `json:"time"`
}

// metrics for disk usage
type DiskMetric struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// metrics for network usage
type NetMetric struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
	ErrIn       uint64 `json:"err_in"`
	ErrOut      uint64 `json:"err_out"`
	DropIn      uint64 `json:"drop_in"`
	DropOut     uint64 `json:"drop_out"`
}

// extended metrics including host details
type MetricMessage struct {
	Host   Host   `json:"host"`
	Metric Metric `json:"metric"`
}
