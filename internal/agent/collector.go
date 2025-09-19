package agent

import (
	"log"
	"monitoring/internal/models"

	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func CollectMetrics() models.Metric {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	cpuInfo, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println("cpu.Percent error:", err)
	}
	cpuPercent := 0.0
	if len(cpuInfo) > 0 {
		cpuPercent = cpuInfo[0]
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
		Hostname: hostname,
		CPU:      cpuPercent,
		RAM:      memPercent,
		Time:     time.Now(),
	}
}
