package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func startStatsd(listenAddr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return fmt.Errorf("invalid UDP address: %v", err)
	}

	l, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("ListenUDP failed: %v", err)
	}

	go listenLoop(l)
	return nil
}

func handleMetric(metricType, metric string, value int64) {
	switch metricType {
	case "c":
		state.IncCounter(metric, value)
	case "ms":
		state.RecordTimer(metric, value)
	case "g":
		state.UpdateGauge(metric, value)
	}
}

func processLine(line string) error {
	line = strings.TrimSpace(line)
	parts := strings.Split(line, "|")
	if len(parts) < 2 {
		return errors.New("metric is missing type")
	}
	parts2 := strings.Split(parts[0], ":")
	metric := parts2[0]
	value, err := strconv.Atoi(parts2[1])
	if err != nil {
		return err
	}

	handleMetric(parts[1], metric, int64(value))
	return nil
}

func listenLoop(l *net.UDPConn) {
	rdr := bufio.NewReader(l)
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("connection error: %v", err)
			}
			return
		}

		if err := processLine(line); err != nil {
			//log.Printf("processLine error: %v", err)
		}
	}
}
