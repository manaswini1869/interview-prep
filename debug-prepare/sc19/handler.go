package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store *DeploymentStore
}

func NewHandler(s *DeploymentStore) *Handler {
	return &Handler{store: s}
}

func (h *Handler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Script string `json:"script"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	d := &Deployment{
		ID:        uuid.New().String(),
		Script:    req.Script,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	h.store.Create(d)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) GetDeployment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	d, ok := h.store.Get(id)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	err := h.store.Delete(id)
	if err != nil {
		if err.Error() == "deployment not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp := map[string]string{"message": "deployment deleted successfully"}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(resp)
}

// BUGGY
func (h *Handler) ListDeployments(w http.ResponseWriter, r *http.Request) {
	deployments := h.store.List()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deployments)
}

// BUGGY
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.store.UpdateStatus(id, req.Status)
	if err != nil {
		if err.Error() == "deployment not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
	resp := map[string]string{"message": "status updated successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
