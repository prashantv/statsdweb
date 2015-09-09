package main

import (
	"sort"
	"sync"
)

const maxTimerSamples = 1000

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
	sort.Sort(Int64Slice(s.Timers[name]))

	// If we have too many samples, drop half of them.
	timerVals := s.Timers[name]
	if numSamples := len(timerVals); numSamples > maxTimerSamples {
		var newTimers []int64
		for i := 0; i < numSamples; i += 2 {
			newTimers = append(newTimers, timerVals[i])
		}
		s.Timers[name] = newTimers
	}
}

// Int64Slice is used to implement Sort for []int64
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (s *metricsState) UpdateGauge(name string, value int64) {
	s.Lock()
	defer s.Unlock()

	s.Gauges[name] = value
}
