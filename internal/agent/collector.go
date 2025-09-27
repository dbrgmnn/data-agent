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

func CollectBaseMetrics() models.BaseMetrics {
	// get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// get host info
	hostInfo, err := host.Info()
	if err != nil {
		log.Println("host.Info error:", err)
	}
	var uptime uint64
	var os, platform, platformVer, kernelVer string
	if hostInfo != nil {
		os = hostInfo.OS
		platform = hostInfo.Platform
		platformVer = hostInfo.PlatformVersion
		kernelVer = hostInfo.KernelVersion
		uptime = hostInfo.Uptime
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

	return models.BaseMetrics{
		Hostname:    hostname,
		OS:          os,
		Platform:    platform,
		PlatformVer: platformVer,
		KernelVer:   kernelVer,
		Uptime:      uptime,
		CPU:         cpuUsage,
		RAM:         memPercent,
		Time:        time.Now(),
	}
}

// disk metrics collection
func CollectDiskMetrics() ([]models.DiskMetrics, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	var metrics []models.DiskMetrics
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		metrics = append(metrics, models.DiskMetrics{
			Path:        p.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}
	return metrics, nil
}

// network metrics collection
func CollectNetMetrics() ([]models.NetMetrics, error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}
	var metrics []models.NetMetrics
	for _, c := range counters {
		metrics = append(metrics, models.NetMetrics{
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
	return metrics, nil
}

// collect all metrics
func CollectAllMetrics(base models.BaseMetrics) models.ExtendedMetrics {
	diskMetrics, _ := CollectDiskMetrics()
	netMetrics, _ := CollectNetMetrics()

	return models.ExtendedMetrics{
		BaseMetrics: base,
		DiskMetrics: diskMetrics,
		NetMetrics:  netMetrics,
	}
}
