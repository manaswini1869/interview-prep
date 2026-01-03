package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type RateLimitRule struct {
	ID     int    `db:"id" json:"id"`
	Method string `db:"method" json:"method"`
	// BUG: Tag mismatch. The DB expects "limit_count" but json is "limit"
	Limit  int    `db:"limit_count" json:"limit"`
	Window string `db:"window" json:"window"`
}

type Server struct {
	DB *sql.DB
	mu sync.Mutex
}

// BUG 2: Broken Pagination
func (s *Server) ListRulesHandler(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(offsetStr)

	// BUG: The SQL syntax below is incorrect for pagination
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id, method, limit_count, window FROM rules LIMIT 10 OFFSET $1", offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rules []RateLimitRule
	for rows.Next() {
		var rule RateLimitRule
		if err := rows.Scan(&rule.ID, &rule.Method, &rule.Limit, &rule.Window); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rules = append(rules, rule)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rules)
}

// TASK: Implement PATCH /rules/:id
// Requirement: If the user only sends {"limit": 100}, the "method" and "window" should remain unchanged.
type UpdateRuleRequest struct {
	Method *string `json:"method"`
	Limit  *int    `json:"limit"`
	Window *string `json:"window"`
}

func (s *Server) PatchRuleHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement partial update logic
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing rule ID", http.StatusBadRequest)
		return
	}
	var req UpdateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.DB.QueryContext(r.Context(), "SELECT method, limit_count, window FROM rules where id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	if !rows.Next() {
		http.Error(w, "Rule not found", http.StatusNotFound)
		return
	}
	var current RateLimitRule
	if err := rows.Scan(&current.Method, &current.Limit, &current.Window); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method != nil {
		current.Method = *req.Method
	}
	if req.Limit != nil {
		current.Limit = *req.Limit
	}
	if req.Window != nil {
		current.Window = *req.Window
	}
	_, err = s.DB.ExecContext(r.Context(), "UPDATE rules SET method = $1, limit_count = $2, window = $3 WHERE id = $4",
		current.Method, current.Limit, current.Window, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
