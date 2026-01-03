package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Deployment struct {
	ID        string    `json:"id"`
	ServiceID string    `json:"service_id"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateDeploymentRequest struct {
	Version string `json:"version"`
}

func (h *Handler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	defer r.Body.Close()

	if serviceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var versionReq CreateDeploymentRequest
	err := json.NewDecoder(r.Body).Decode(&versionReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if versionReq.Version == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d := Deployment{
		ID:        uuid.New().String(),
		ServiceID: serviceID,
		Version:   versionReq.Version,
		CreatedAt: time.Now(),
	}

	err = h.store.CreateDeployment(r.Context(), &d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetServiceConfig(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")

	if serviceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cfg, err := h.store.GetConfig(r.Context(), serviceID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PostRollbackDeployment(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	deploymentID := chi.URLParam(r, "deployment_id")

	if serviceID == "" || deploymentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.store.RollbackDeployment(r.Context(), serviceID, deploymentID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "rollback successful"})
	w.WriteHeader(http.StatusNoContent)
}
