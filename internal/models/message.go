package models

// extended metrics including host details
type MetricMessage struct {
	Host   Host   `json:"host"`
	Metric Metric `json:"metric"`
}

// constructor for MetricMessage
func NewMetricMessage(host *Host, metric *Metric) *MetricMessage {
	return &MetricMessage{
		Host:   *host,
		Metric: *metric,
	}
}
