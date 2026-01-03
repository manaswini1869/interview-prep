package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type ServiceConfig struct {
	ServiceID string            `json:"service_id"`
	Settings  map[string]string `json:"settings"`
	Version   int64             `json:"version"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type UpdateConfigRequest struct {
	Settings map[string]string `json:"settings"`
	Version  int64             `json:"version"`
}

func (h *Handler) UpdateServiceConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req UpdateConfigRequest
	// var cfg ServiceConfig
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	version := req.Version
	settings := req.Settings
	if version <= 0 || settings == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cfg ServiceConfig
	cfg.Version = version
	cfg.Settings = settings
	cfg.ServiceID = serviceID
	cfg.UpdatedAt = time.Now()

	updated, err = h.store.UpdateConfigifVersionMatches(r.Context(), &cfg)
	if err != nil {
		if err == store.ErrVersionConflict {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !updated {
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
	w.WriteHeader(http.StatusOK)
}
