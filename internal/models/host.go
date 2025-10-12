package models

// represents a monitored host system
type Host struct {
	ID          int64  `json:"id"`
	Hostname    string `json:"hostname"`
	OS          string `json:"os"`
	Platform    string `json:"platform"`
	PlatformVer string `json:"platformver"`
	KernelVer   string `json:"kernelver"`
}

// constructor for Host with basic validation
func NewHost(hostname, os, platform, platformVer, kernelVer string) (*Host, error) {
	if hostname == "" {
		hostname = "unknown"
	}

	return &Host{
		Hostname:    hostname,
		OS:          os,
		Platform:    platform,
		PlatformVer: platformVer,
		KernelVer:   kernelVer,
	}, nil
}
