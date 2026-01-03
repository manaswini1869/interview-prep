package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Deployment struct {
	ID        string `json:"id"`
	ServiceID string `json:"service_id"`
	Version   string `json:"version"`
	CreatedAt string `json:"created_at"`
}

type ListDeploymentsResponse struct {
	Deployments []Deployment `json:"deployments"`
	TotalCount  int          `json:"total_count"`
	NextCursor  string       `json:"next_offset,omitempty"`
	HasMore     bool         `json:"has_more,omitempty"`
}

const (
	DefaultListLimit = 20
	MaxListLimit     = 100
)

func (h *Handler) ListDeployments(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	limit := DefaultListLimit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limitParsed, err := strconv.Atoi(limitStr)
		if err != nil || limitParsed <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		limit = limitParsed
		if limit > MaxListLimit {
			limit = MaxListLimit
		}
	}
	cursor, _ := strconv.Atoi(r.URL.Query().Get("cursor"))
	var deployments []Deployment
	deployments, err := h.store.ListDeployments(
		r.Context(), serviceID, limit+1, cursor,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hasMore := len(deployments) > limit
	if hasMore {
		deployments = deployments[:limit]
	}
	var nextCursor string
	if hasMore && len(deployments) > 0 {
		nextCursor = deployments[len(deployments)-1].ID // Assuming ID can be converted to int
	}
	totalCount, err := h.store.CountDeployments(r.Context(), serviceID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ListDeploymentsResponse := ListDeploymentsResponse{
		Deployments: deployments,
		TotalCount:  totalCount,
		HasMore:     hasMore,
		NextCursor:  nextCursor,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListDeploymentsResponse)
}
