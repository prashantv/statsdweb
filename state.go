package main

import "sync"

type metricsState struct {
	sync.RWMutex

	Counters map[string]int64   `json:"counters"`
	Gauges   map[string]int64   `json:"gauges"`
	Timers   map[string][]int64 `json:"timers"`
}

func newMetricsState() *metricsState {
	return &metricsState{
		Counters: make(map[string]int64),
		Timers:   make(map[string][]int64),
		Gauges:   make(map[string]int64),
	}
}

func (s *metricsState) IncCounter(name string, value int64) {
	s.Lock()
	defer s.Unlock()

	s.Counters[name] += value
}

func (s *metricsState) RecordTimer(name string, value int64) {
	s.Lock()
	defer s.Unlock()

	s.Timers[name] = append(s.Timers[name], value)
}

func (s *metricsState) UpdateGauge(name string, value int64) {
	s.Lock()
	defer s.Unlock()

	s.Gauges[name] = value
}
