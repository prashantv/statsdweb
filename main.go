package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sort"
)

var (
	listenAddr = flag.String("listenAddr", "127.0.0.1:8125", "Host:Port to listen on")

	state = newMetricsState()
)

func init() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/counters", countersHandler)
}

func main() {
	flag.Parse()
	if err := startStatsd(*listenAddr); err != nil {
		log.Fatalf("statsd failed: %v", err)
	}
	log.Printf("Start listening on :9090")
	log.Fatalf("http.ListenAndServe failed: %v", http.ListenAndServe(":9090", nil))
}

type record struct {
	CounterName string `json:"counterName"`
	Value       int64  `json:"value"`
}

type recordByName []record

func (s recordByName) Len() int           { return len(s) }
func (s recordByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s recordByName) Less(i, j int) bool { return s[i].CounterName < s[j].CounterName }

func countersHandler(w http.ResponseWriter, r *http.Request) {
	state.RLock()
	state.RUnlock()

	var records []record
	for k, v := range state.Counters {
		records = append(records, record{k, v})
	}
	sort.Sort(recordByName(records))

	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(records)
}
