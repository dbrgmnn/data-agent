package models

import "time"

type Metric struct {
	Hostname string    `json:"hostname"`
	CPU      float64   `json:"cpu"`
	RAM      float64   `json:"ram"`
	Time     time.Time `json:"time"`
}
