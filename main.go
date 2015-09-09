package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var (
	listenAddr = flag.String("listenAddr", "127.0.0.1:8125", "Host:Port to listen on")

	state = newMetricsState()
)

func init() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/state", stateHandler)
}

func main() {
	flag.Parse()
	if err := startStatsd(*listenAddr); err != nil {
		log.Fatalf("statsd failed: %v", err)
	}
	log.Printf("Start listening on :9090")
	log.Fatalf("http.ListenAndServe failed: %v", http.ListenAndServe(":9090", nil))
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	state.RLock()
	state.RUnlock()

	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(state)
}
