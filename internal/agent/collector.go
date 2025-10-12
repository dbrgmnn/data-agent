package agent

import (
	"log"
	"monitoring/internal/models"

	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// collect host information
func CollectHostInfo() models.Host {
	// get hostname
	hostname, _ := os.Hostname()

	// get host info
	hostInfo, err := host.Info()
	if err != nil {
		log.Println("host.Info error:", err)
	}
	if hostInfo == nil {
		log.Println("host.Info returned nil")
		return models.Host{}
	}

	// create Host model
	host, err := models.NewHost(
		hostname,
		hostInfo.OS,
		hostInfo.Platform,
		hostInfo.PlatformVersion,
		hostInfo.KernelVersion,
	)
	if err != nil {
		log.Println("NewHost error:", err)
	}

	return *host
}

// collect metric information
func CollectMetricInfo() models.Metric {
	// get uptime
	uptime, err := host.Uptime()
	if err != nil {
		log.Println("host.Uptime error:", err)
	}

	// get CPU info
	cpuInfo, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println("cpu.Percent error:", err)
	}
	cpuUsage := 0.0
	if len(cpuInfo) > 0 {
		cpuUsage = cpuInfo[0]
	}

	// get memory info
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("mem.VirtualMemory error:", err)
	}
	memPercent := 0.0
	if memInfo != nil {
		memPercent = memInfo.UsedPercent
	}

	diskMetric, err := CollectDiskMetric()
	if err != nil {
		log.Println("CollectDiskMetric error:", err)
	}
	netMetric, err := CollectNetMetric()
	if err != nil {
		log.Println("CollectNetMetric error:", err)
	}

	// create Metric model
	metric, err := models.NewMetric(
		uptime,
		cpuUsage,
		memPercent,
		diskMetric,
		netMetric,
	)
	if err != nil {
		log.Println("NewMetric error:", err)
	}

	return *metric
}

// collect disk metrics
func CollectDiskMetric() ([]models.DiskMetric, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	var metrics []models.DiskMetric
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		metrics = append(metrics, models.DiskMetric{
			Path:        p.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	if metrics == nil {
		metrics = []models.DiskMetric{}
	}
	return metrics, nil
}

// collect network metrics
func CollectNetMetric() ([]models.NetMetric, error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}
	var metrics []models.NetMetric
	for _, c := range counters {
		metrics = append(metrics, models.NetMetric{
			Name:        c.Name,
			BytesSent:   c.BytesSent,
			BytesRecv:   c.BytesRecv,
			PacketsSent: c.PacketsSent,
			PacketsRecv: c.PacketsRecv,
			ErrIn:       c.Errin,
			ErrOut:      c.Errout,
			DropIn:      c.Dropin,
			DropOut:     c.Dropout,
		})
	}

	if metrics == nil {
		metrics = []models.NetMetric{}
	}
	return metrics, nil
}
