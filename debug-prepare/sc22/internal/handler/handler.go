package handler

import (
	"encoding/json"
	"net/http"

	"sc22/internal/routes"
)

type Handler struct {
	Store *routes.Store
}

func (h *Handler) ListRoutes(w http.ResponseWriter, r *http.Request) {
	routes := h.Store.List()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes) // missing status code
}

func (h *Handler) EnableRoute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // incorrect source
	if !h.Store.Enable(id) {
		http.Error(w, "not found", http.StatusNotFound) // wrong status
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
