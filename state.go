package main

import "sync"

type metricsState struct {
	sync.RWMutex

	Counters map[string]int64 `json:"counters"`
}

func newMetricsState() *metricsState {
	return &metricsState{
		Counters: make(map[string]int64),
	}
}

func (s *metricsState) IncCounter(name string, value int64) {
	s.Lock()
	defer s.Unlock()

	s.Counters[name] += value
}
