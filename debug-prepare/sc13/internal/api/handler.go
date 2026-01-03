package api

import (
	"encoding/json"
	"net/http"

	"workers/internal/model"
	"workers/internal/store"
)

type Handler struct {
	Store *store.MemoryStore
}

func (h *Handler) CreateWorker(w http.ResponseWriter, r *http.Request) {
	var worker model.Worker
	err := json.NewDecoder(r.Body).Decode(&worker)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if worker.ID == "" || worker.Script == "" || worker.Version <= 0 || worker.Active == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// worker.Version = 1
	// worker.Active = true

	err = h.Store.Create(worker)
	if err != nil {
		if err.Error() == "worker already exists" {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(worker)
}

func (h *Handler) GetWorker(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	worker, err := h.Store.Get(id)
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(worker)
}

// TODO: Implement DeactivateWorker
func (h *Handler) DeactivateWorker(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	worker, err := h.Store.DeactivateWorker(id)
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(worker)
}
