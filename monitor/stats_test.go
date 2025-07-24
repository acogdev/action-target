package monitor

import (
	"testing"
	"time"
)

func TestAddResult(t *testing.T) {
	stats := &Stats{host: "example.com"}

	// Test successful result
	duration := 100 * time.Millisecond
	stats.addResult(true, duration)

	if stats.Sent != 1 {
		t.Errorf("Expected Sent to be 1, got %d", stats.Sent)
	}
	if stats.Received != 1 {
		t.Errorf("Expected Received to be 1, got %d", stats.Received)
	}
	if !stats.Up {
		t.Error("Expected Up to be true")
	}
	if stats.MinTime != duration {
		t.Errorf("Expected MinTime to be %v, got %v", duration, stats.MinTime)
	}
	if stats.MaxTime != duration {
		t.Errorf("Expected MaxTime to be %v, got %v", duration, stats.MaxTime)
	}
	if stats.totalTime != duration {
		t.Errorf("Expected totalTime to be %v, got %v", duration, stats.totalTime)
	}
	if len(stats.times) != 1 || stats.times[0] != duration {
		t.Errorf("Expected times to contain [%v], got %v", duration, stats.times)
	}

	// Test failed result
	stats.addResult(false, 200*time.Millisecond)

	if stats.Sent != 2 {
		t.Errorf("Expected Sent to be 2, got %d", stats.Sent)
	}
	if stats.Received != 1 {
		t.Errorf("Expected Received to still be 1, got %d", stats.Received)
	}
	if stats.Up {
		t.Error("Expected Up to be false after failed result")
	}
	// Verify that failed results don't affect timing stats
	if stats.MinTime != duration {
		t.Errorf("Expected MinTime to remain %v, got %v", duration, stats.MinTime)
	}
	if stats.MaxTime != duration {
		t.Errorf("Expected MaxTime to remain %v, got %v", duration, stats.MaxTime)
	}
	if len(stats.times) != 1 {
		t.Errorf("Expected times length to remain 1, got %d", len(stats.times))
	}
}

func TestAddResult_MinMaxTracking(t *testing.T) {
	stats := &Stats{host: "example.com"}

	durations := []time.Duration{
		150 * time.Millisecond,
		50 * time.Millisecond,  // min
		300 * time.Millisecond, // max
		100 * time.Millisecond,
	}

	for _, d := range durations {
		stats.addResult(true, d)
	}

	expectedMin := 50 * time.Millisecond
	expectedMax := 300 * time.Millisecond

	if stats.MinTime != expectedMin {
		t.Errorf("Expected MinTime to be %v, got %v", expectedMin, stats.MinTime)
	}
	if stats.MaxTime != expectedMax {
		t.Errorf("Expected MaxTime to be %v, got %v", expectedMax, stats.MaxTime)
	}
}

func TestGetAverage(t *testing.T) {
	stats := &Stats{host: "example.com"}

	// Test with no received packets
	avg := stats.GetAverage()
	if avg != 0 {
		t.Errorf("Expected average to be 0 with no received packets, got %v", avg)
	}

	durations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		300 * time.Millisecond,
	}

	for _, d := range durations {
		stats.addResult(true, d)
	}

	expectedAvg := 200 * time.Millisecond // (100+200+300)/3
	avg = stats.GetAverage()
	if avg != expectedAvg {
		t.Errorf("Expected average to be %v, got %v", expectedAvg, avg)
	}

	// Failed result should not affect average
	stats.addResult(false, 1000*time.Millisecond)
	avg = stats.GetAverage()
	if avg != expectedAvg {
		t.Errorf("Expected average to remain %v after failed result, got %v", expectedAvg, avg)
	}
}


func TestGetPacketLoss(t *testing.T) {
	stats := &Stats{host: "example.com"}

	// Test with no packets sent
	loss := stats.GetPacketLoss()
	if loss != 0 {
		t.Errorf("Expected packet loss to be 0 with no packets sent, got %f", loss)
	}

	// Test with all packets successful
	stats.addResult(true, 100*time.Millisecond)
	stats.addResult(true, 200*time.Millisecond)
	loss = stats.GetPacketLoss()
	if loss != 0 {
		t.Errorf("Expected packet loss to be 0 with all packets successful, got %f", loss)
	}

	// Test with some packet loss
	stats.addResult(false, 100*time.Millisecond)
	stats.addResult(false, 100*time.Millisecond)
	loss = stats.GetPacketLoss()
	expected := 50.0
	if loss != expected {
		t.Errorf("Expected packet loss to be %f, got %f", expected, loss)
	}

	// Test with 100% packet loss
	statsAllFailed := &Stats{host: "example.com"}
	statsAllFailed.addResult(false, 100*time.Millisecond)
	statsAllFailed.addResult(false, 200*time.Millisecond)
	loss = statsAllFailed.GetPacketLoss()
	expected = 100.0
	if loss != expected {
		t.Errorf("Expected packet loss to be %f with all packets failed, got %f", expected, loss)
	}
}
