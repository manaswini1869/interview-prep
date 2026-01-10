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
	json.NewDecoder(r.Body).Decode(&route)

	h.Store.Save(route)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetRoute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	route, ok := h.Store.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(route)
}

// TODO: Implement
func (h *Handler) GetConflicts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
