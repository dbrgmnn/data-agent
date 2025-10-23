package models

import (
	"errors"
	"log"
	"time"
)

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

// constructor for Metric
func NewMetric(uptime uint64, cpu, ram float64, disk []DiskMetric, network []NetMetric) *Metric {

	metric := &Metric{
		Uptime:  uptime,
		CPU:     cpu,
		RAM:     ram,
		Disk:    disk,
		Network: network,
		Time:    time.Now(),
	}
	if err := metric.Validate(); err != nil {
		log.Println("Metric validation error:", err)
		return nil
	}

	return metric
}

// validation metrics
func (m *Metric) Validate() error {
	if m.Uptime == 0 {
		return errors.New("invalid uptime value")
	}
	if m.CPU < 0 || m.CPU > 100 {
		return errors.New("invalid cpu value")
	}
	if m.RAM < 0 || m.RAM > 100 {
		return errors.New("invalid ram value")
	}
	return nil
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
