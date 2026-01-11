package api

import (
	"encoding/json"
	"net/http"

	"sc23/internal/model"
	"sc23/internal/store"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Store *store.MemoryStore
}

func (h *Handler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var route model.Route
	err := json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	modelRoute, ok := h.Store.Save(route)
	if !ok {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(modelRoute)
}

func (h *Handler) GetRoute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	route, ok := h.Store.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(route)
}

// TODO: Implement
func (h *Handler) GetConflicts(w http.ResponseWriter, r *http.Request) {
	conflicts := h.Store.GetConflicts()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conflicts)
}
