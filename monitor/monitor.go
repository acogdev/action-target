package monitor

import (
	"fmt"
	"net"
	"strconv"
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
}

func Monitor(hosts []string, port string, interval int, configFile string) {
	timeout := 2 * time.Second

	if len(configFile) > 1 {
		fmt.Printf("Using values from %s\n", configFile)
		config := ReadConfig(configFile)
		hosts = config.Hosts
		port = strconv.Itoa(config.Port)
		interval = config.Interval
	}

	fmt.Printf("The following services will checked every %v seconds on port %s: %s ", interval, port, hosts)

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
