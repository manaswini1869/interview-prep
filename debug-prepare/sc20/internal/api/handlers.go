package api

import (
	"encoding/json"
	"net/http"

	"github.com/manaswini1869/debug-prepare/sc20/internal/store"

	"github.com/manaswini1869/debug-prepare/sc20/internal/model"
)

type Handler struct {
	Store *store.RolloutStore
}

func (h *Handler) CreateRollout(w http.ResponseWriter, r *http.Request) {
	var rollout model.Rollout
	err := json.NewDecoder(r.Body).Decode(&rollout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rollout.Completed = false
	err = h.Store.Save(&rollout)
	if err != nil {
		if err.Error() == "duplicate ID" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := map[string]string{"id": rollout.ID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetRollout(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	rollout, err := h.Store.Get(id)
	if err != nil {
		if err.Error() == "rollout not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rollout)
}

// TODO: Implement UpdateRegionStatus
func (h *Handler) UpdateRegionStatus(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	var statusUpdate model.StatusUpdate
	err := json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Store.UpdateRegionStatus(id, statusUpdate)
	if err != nil {
		if err.Error() == "rollout not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
