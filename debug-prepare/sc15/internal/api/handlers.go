package api

import (
	"encoding/json"
	"net/http"

	"token/internal/model"
	"token/internal/store"
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

	h.Store.Save(t)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing token ID", http.StatusBadRequest)
		return
	}
	token, err := h.Store.Get(id)
	if err != nil {
		http.Error(w, "Token not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

// TODO: Implement RevokeToken
func (h *Handler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing token ID", http.StatusBadRequest)
		return
	}
	err := h.Store.Revoke(id)
	if err != nil {
		if err.Error() == "token not found" {
			http.Error(w, "Token not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(`"message": "Token revoked successfully"`))
}
