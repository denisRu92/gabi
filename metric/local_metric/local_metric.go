package local_metric

import (
	"palo-alto/metric"
	"sync/atomic"
	"time"
)

type localMetric struct {
	wordCounter             int64
	requestCounter          int64
	requestProcessingTimeNs []int64

	avgProcessingTimeCh chan avgProcessingTimeReq
	addReqProcessTimeCh chan int64
	stopCh              chan struct{}
}

type avgProcessingTimeReq struct {
	respCh chan int64
}

// New return local_metric
func New() metric.Metric {
	return &localMetric{
		wordCounter:             0,
		requestCounter:          0,
		requestProcessingTimeNs: []int64{},

		avgProcessingTimeCh: make(chan avgProcessingTimeReq),
		addReqProcessTimeCh: make(chan int64),
		stopCh:              make(chan struct{}),
	}
}

func (m *localMetric) Start() {
	for {
		select {
		case duration := <-m.addReqProcessTimeCh:
			m.addProcessingTiming(duration)
		case req := <-m.avgProcessingTimeCh:
			req.respCh <- m.avgProcessingTimeNs()
		case <-m.stopCh:
			return
		}
	}
}

func (m *localMetric) Stop() {
	close(m.stopCh)
}

func (m *localMetric) GetStates() metric.States {
	respCh := make(chan int64)
	m.avgProcessingTimeCh <- avgProcessingTimeReq{
		respCh: respCh,
	}
	avgProcessingTimeNs := <-respCh

	return metric.States{
		TotalWords:          m.wordCounter,
		TotalRequests:       m.requestCounter,
		AvgProcessingTimeNs: avgProcessingTimeNs,
	}
}

func (m *localMetric) IncWordCounter() {
	atomic.AddInt64(&m.wordCounter, 1)
}

func (m *localMetric) IncRequestCounter() {
	atomic.AddInt64(&m.requestCounter, 1)
}

func (m *localMetric) AddProcessingTiming(start time.Time) {
	m.addReqProcessTimeCh <- time.Since(start).Nanoseconds()
}

func (m *localMetric) addProcessingTiming(duration int64) {
	m.requestProcessingTimeNs = append(m.requestProcessingTimeNs, duration)
}

func (m *localMetric) avgProcessingTimeNs() int64 {
	var avgProcessingTimeNs int64
	if len(m.requestProcessingTimeNs) > 0 {
		for _, duration := range m.requestProcessingTimeNs {
			avgProcessingTimeNs += duration
		}

		avgProcessingTimeNs /= int64(len(m.requestProcessingTimeNs))
	}

	return avgProcessingTimeNs
}
