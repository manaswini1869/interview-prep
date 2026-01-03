package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/sync/errgroup"
)

var (
	opsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "config_updates_total",
		Help: "The total number of configuration updates",
	})
)

type ConfigEntry struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Server struct {
	DB      *sql.DB
	Regions []string
	mu      sync.Mutex
}

func (s *Server) UpdateConfigHandler(w http.ResponseWriter, r *http.Request) {
	var entry ConfigEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	_, err := s.DB.ExecContext(r.Context(), "INSERT INTO configs (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = $2", entry.Key, entry.Value)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}
	opsCounter.Inc()
	s.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

// TASK: Implement this endpoint.
// It should iterate through s.Regions and call s.pingRegion(ctx, region) concurrently.
func (s *Server) SyncConfigHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation using errgroup
	ctx := r.Context()
	g, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex
	var failedRegions []string

	for _, region := range s.Regions {
		region := region // capture range variable
		g.Go(func() error {
			if err := s.pingRegion(ctx, region); err != nil {
				mu.Lock()
				failedRegions = append(failedRegions, region)
				mu.Unlock()
				return err
			}
			return nil
		})
	}
	_ = g.Wait()
	if len(failedRegions) > 0 {
		http.Error(w, "Failed regions: "+string(len(failedRegions)), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All regions synced successfully"))

}

func (s *Server) pingRegion(ctx context.Context, region string) error {
	// Simulates an external API call
	req, err := http.NewRequestWithContext(ctx, "POST", "https://"+region+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
