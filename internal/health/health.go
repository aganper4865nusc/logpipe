// Package health provides a simple HTTP health-check endpoint for logpipe.
package health

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the current health state of the process.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	StartedAt time.Time `json:"started_at"`
}

// Server exposes a /healthz HTTP endpoint.
type Server struct {
	addr      string
	startedAt time.Time
	ready     atomic.Bool
	server    *http.Server
}

// New creates a new health Server listening on addr (e.g. ":8080").
func New(addr string) *Server {
	s := &Server{
		addr:      addr,
		startedAt: time.Now(),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return s
}

// SetReady marks the server as ready to serve traffic.
func (s *Server) SetReady(ready bool) {
	s.ready.Store(ready)
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Close shuts down the HTTP server gracefully.
func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ok := s.ready.Load()
	status := Status{
		OK:        ok,
		Uptime:    time.Since(s.startedAt).Round(time.Second).String(),
		StartedAt: s.startedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	_ = json.NewEncoder(w).Encode(status)
}
