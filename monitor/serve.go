package monitor

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
	"embed"
)

//go:embed templates/*
var content embed.FS

func (s *Serve) renderStats(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(content, "templates/layout.html"))
	tmpl.Execute(w, s.CurrentStats)
}

type JsonStats struct {
	Host            string        `json:"host"`
	Up              bool          `json:"up"`
	Sent            int           `json:"sent"`
	Received        int           `json:"received"`
	PacketLoss      float64       `json:"packetLoss"`
	MinTime         time.Duration `json:"minTime"`
	MaxTime         time.Duration `json:"maxTime"`
	AverageResponse time.Duration `json:"average"`
}

func (s *Serve) Stats(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	data := []JsonStats{}

	for _, v := range s.CurrentStats {
		d := JsonStats{
			Host:            v.Host,
			Sent:            v.Sent,
			Received:        v.Received,
			MinTime:         time.Duration(v.MinTime.Milliseconds()),
			MaxTime:         time.Duration(v.MaxTime.Milliseconds()),
			Up:              v.Up,
			PacketLoss:      v.GetPacketLoss(),
			AverageResponse: time.Duration(v.GetAverage().Milliseconds()),
		}
		data = append(data, d)
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func (s *Serve) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("up"))
}

func (s *Serve) Serve() {
	http.HandleFunc("/", s.renderStats)
	http.HandleFunc("/stats", s.Stats)
	http.HandleFunc("/health", s.health)
	http.ListenAndServe(":8080", nil)
}

type Serve struct {
	CurrentStats map[string]*Stats
}
