package models

// extended metrics including host details
type MetricMessage struct {
	Host   Host   `json:"host"`
	Metric Metric `json:"metric"`
}
