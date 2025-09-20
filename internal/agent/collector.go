package agent

import (
	"log"
	"monitoring/internal/models"

	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func CollectMetrics() models.Metric {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	hostInfo, err := host.Info()
	if err != nil {
		log.Println("host.Info error:", err)
	}

	cpuInfo, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println("cpu.Percent error:", err)
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("mem.VirtualMemory error:", err)
	}

	memPercent := 0.0
	if memInfo != nil {
		memPercent = memInfo.UsedPercent
	}

	return models.Metric{
		Hostname:    hostname,
		OS:          hostInfo.OS,
		Platform:    hostInfo.Platform,
		PlatformVer: hostInfo.PlatformVersion,
		KernelVer:   hostInfo.KernelVersion,
		Uptime:      hostInfo.Uptime,
		CPU:         cpuInfo[0],
		RAM:         memPercent,
		Time:        time.Now(),
	}
}
