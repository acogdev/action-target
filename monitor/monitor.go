package monitor

import (
	"fmt"
	"net"
	"time"
)

func isHostUp(host string, port string, timeout time.Duration, stat *Stats) {
	address := net.JoinHostPort(host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	duration := time.Since(start)
	success := false
	if err == nil {
		success = true
		defer conn.Close()
	}
	stat.addResult(success, duration)
	printSummary(stat)
}

func Monitor(hosts []string, port string, interval int) {
	fmt.Println("called the monitor code with ", hosts, port, interval)
	timeout := 2 * time.Second

	var currentStats = make(map[string]*Stats)

	for _, host := range hosts {
		currentStats[host] = &Stats{Host: host}
	}

	server := &Serve{CurrentStats: currentStats}
	go server.Serve()

	for {
		for _, host := range hosts {
			go isHostUp(host, port, timeout, currentStats[host])
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}

}
