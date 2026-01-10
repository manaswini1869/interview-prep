package api

import (
	"encoding/json"
	"net/http"

	"sc21/internal/model"
	"sc21/internal/store"
)

type Handler struct {
	Store *store.Store
}

func (h *Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	var t model.Token
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if t.ID == "" {
		http.Error(w, "Token ID is required", http.StatusBadRequest)
		return
	}

	t, err = h.Store.Save(t)
	if err != nil {
		if err.Error() == "token already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Token ID is required", http.StatusBadRequest)
		return
	}
	token, err := h.Store.Get(id)
	if err != nil {
		if err.Error() == "token not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

// TODO: Implement RevokeToken
func (h *Handler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Token ID is required", http.StatusBadRequest)
		return
	}
	token, err := h.Store.Revoke(id)
	if err != nil {
		if err.Error() == "token not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
