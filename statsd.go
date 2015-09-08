package main

import (
	"bufio"
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

func processLine(line string) error {
	line = strings.TrimSpace(line)
	parts := strings.Split(line, "|")
	if parts[1] != "c" {
		return fmt.Errorf("ignoring type: %v", parts[1])
	}
	parts2 := strings.Split(parts[0], ":")
	metric := parts2[0]
	count, err := strconv.Atoi(parts2[1])
	if err != nil {
		return err
	}

	state.IncCounter(metric, int64(count))
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
