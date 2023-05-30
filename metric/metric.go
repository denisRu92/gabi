package metric

import "time"

type Metric interface {
	Start()
	Stop()
	AddProcessingTiming(start time.Time)
	IncWordCounter()
	IncRequestCounter()
	GetStates() States
}

type States struct {
	TotalWords          int64 `json:"totalWords"`
	TotalRequests       int64 `json:"totalRequests"`
	AvgProcessingTimeNs int64 `json:"avgProcessingTimeNs"`
}
