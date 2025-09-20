package models

import "time"

type Metric struct {
	Hostname    string    `json:"hostname"`
	OS          string    `json:"os"`
	Platform    string    `json:"platform"`
	PlatformVer string    `json:"platformver"`
	KernelVer   string    `json:"kernelver"`
	Uptime      uint64    `json:"uptime"`
	CPU         float64   `json:"cpu"`
	RAM         float64   `json:"ram"`
	Time        time.Time `json:"time"`
}
