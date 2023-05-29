package metric

import (
	"sync/atomic"
	"time"
)

type Metric struct {
	wordCounter                int64
	requestCounter             int64
	requestProcessingTimeNsSum int64
	processedRequestCounter    int64
}

type States struct {
	TotalWords          int64 `json:"totalWords"`
	TotalRequests       int64 `json:"totalRequests"`
	AvgProcessingTimeNs int64 `json:"avgProcessingTimeNs"`
}

// New return Metric
func New() *Metric {
	return &Metric{
		wordCounter:                0,
		requestCounter:             0,
		requestProcessingTimeNsSum: 0,
		processedRequestCounter:    0,
	}
}

func (m *Metric) IncWordCounter() {
	atomic.AddInt64(&m.wordCounter, 1)
}

func (m *Metric) AddWordCounter(count int64) {
	atomic.AddInt64(&m.wordCounter, count)
}

func (m *Metric) IncRequestCounter() {
	atomic.AddInt64(&m.requestCounter, 1)
}

func (m *Metric) ObserveProcessingTiming(start time.Time) {
	atomic.AddInt64(&m.processedRequestCounter, 1)
	atomic.AddInt64(&m.requestProcessingTimeNsSum, time.Since(start).Nanoseconds())
}

func (m *Metric) GetStates() States {
	var avgProcessingTimeNs int64
	if m.processedRequestCounter > 0 {
		avgProcessingTimeNs = m.requestProcessingTimeNsSum / m.processedRequestCounter
	}

	return States{
		TotalWords:          m.wordCounter,
		TotalRequests:       m.requestCounter,
		AvgProcessingTimeNs: avgProcessingTimeNs,
	}
}
