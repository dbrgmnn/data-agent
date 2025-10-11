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
