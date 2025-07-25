package monitor

import (
	"time"
)

type Stats struct {
	Host      string        
	Sent      int           
	Received  int           
	MinTime   time.Duration 
	MaxTime   time.Duration 
	totalTime time.Duration
	times     []time.Duration
	Up        bool          
}

func (s *Stats) addResult(success bool, duration time.Duration) {
	s.Sent++
	s.Up = success
	if success {
		s.Received++
		s.totalTime += duration
		s.times = append(s.times, duration)

		if s.MinTime == 0 || duration < s.MinTime {
			s.MinTime = duration
		}
		if duration > s.MaxTime {
			s.MaxTime = duration
		}
	}
}

func (s *Stats) GetAverage() time.Duration {
	if s.Received == 0 {
		return 0
	}
	return s.totalTime / time.Duration(s.Received)
}

func (s *Stats) GetPacketLoss() float64 {
	if s.Sent == 0 {
		return 0
	}
	return float64(s.Sent-s.Received) / float64(s.Sent) * 100
}
