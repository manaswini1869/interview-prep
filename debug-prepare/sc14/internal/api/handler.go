package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"rollout/internal/model"
	"rollout/internal/store"
)

type Handler struct {
	Store *store.RolloutStore
}

func (h *Handler) CreateRollout(w http.ResponseWriter, r *http.Request) {
	var rollout model.Rollout
	err := json.NewDecoder(r.Body).Decode(&rollout)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Store.Save(&rollout)
	if err != nil {
		if err.Error() == "rollout found" {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(rollout.ID))
}

func (h *Handler) GetRollout(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rollout, err := h.Store.Get(id)

	if err != nil {
		if err.Error() == "rollout not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rollout)
}

// TODO: Implement UpdateRegionStatus
func (h *Handler) DeactivateRollout(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.Store.Deactivate(id)

	if err != nil {
		if err.Error() == "rollout not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	responseMessage := fmt.Sprintf("Updated the status of %s", id)
	w.Write([]byte(responseMessage))
}
