package main

import (
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type PingStats struct {
	sent      int
	received  int
	minTime   time.Duration
	maxTime   time.Duration
	totalTime time.Duration
	times     []time.Duration
}

func (s *PingStats) addResult(success bool, duration time.Duration) {
	s.sent++
	if success {
		s.received++
		s.totalTime += duration
		s.times = append(s.times, duration)

		if s.minTime == 0 || duration < s.minTime {
			s.minTime = duration
		}
		if duration > s.maxTime {
			s.maxTime = duration
		}
	}
}

func (s *PingStats) getAverage() time.Duration {
	if s.received == 0 {
		return 0
	}
	return s.totalTime / time.Duration(s.received)
}

func (s *PingStats) getStdDev() time.Duration {
	if len(s.times) <= 1 {
		return 0
	}

	avg := s.getAverage()
	var sum float64

	for _, t := range s.times {
		diff := float64(t - avg)
		sum += diff * diff
	}

	variance := sum / float64(len(s.times)-1)
	stddev := time.Duration(math.Sqrt(variance))
	return stddev
}

func (s *PingStats) getPacketLoss() float64 {
	if s.sent == 0 {
		return 0
	}
	return float64(s.sent-s.received) / float64(s.sent) * 100
}

func tcpPing(host string, port string, timeout time.Duration) (bool, time.Duration) {
	address := net.JoinHostPort(host, port)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", address, timeout)
	duration := time.Since(start)

	if err != nil {
		return false, duration
	}
	defer conn.Close()
	return true, duration
}

func sample() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <host> <port>\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]
	port := os.Args[2]
	timeout := 3 * time.Second
	interval := 1 * time.Second

	stats := &PingStats{}

	// Handle Ctrl+C gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		printSummary(host, port, stats)
		os.Exit(0)
	}()

	fmt.Printf("TCP PING %s:%s\n", host, port)

	seq := 1
	for {
		success, duration := tcpPing(host, port, timeout)
		stats.addResult(success, duration)

		if success {
			fmt.Printf("Connected to %s:%s: seq=%d time=%v\n",
				host, port, seq, duration.Truncate(time.Microsecond))
		} else {
			fmt.Printf("Connection to %s:%s failed: seq=%d time=%v\n",
				host, port, seq, duration.Truncate(time.Microsecond))
		}

		seq++
		time.Sleep(interval)
	}
}

func printSummary(host string, port string, stats *PingStats) {
	fmt.Printf("\n--- %s:%s ping statistics ---\n", host, port)
	fmt.Printf("%d packets transmitted, %d received, %.1f%% packet loss\n",
		stats.sent, stats.received, stats.getPacketLoss())

	if stats.received > 0 {
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.minTime.Truncate(time.Microsecond),
			stats.getAverage().Truncate(time.Microsecond),
			stats.maxTime.Truncate(time.Microsecond),
			stats.getStdDev().Truncate(time.Microsecond))
	}
}
