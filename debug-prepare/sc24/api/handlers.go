package api

import (
	"encoding/json"
	"net/http"
	"sc24/models"
	"sc24/store"

	"github.com/gorilla/mux"
)

type Server struct {
	Store *store.Store
}

// GetWorkerHandler handles GET /workers/{name}
func (s *Server) GetWorkerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "Missing worker name", http.StatusBadRequest)
		return
	}

	worker, err := s.Store.GetWorkerByName(name)
	if err != nil {
		// BUG 1 (Consumer side): It sends 500 Internal Server Error for everything,
		// even if the user just queried a non-existent worker (should be 404).
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Worker Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(worker)
}

// UpdateCPUHandler handles PUT /workers/{name}/cpu
func (s *Server) UpdateCPUHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		http.Error(w, "Missing worker name", http.StatusBadRequest)
		return
	}

	type Request struct {
		Limit int `json:"limit"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	// BUG 2 (Consumer side): The store returns nil error even if row wasn't found (0 rows affected).
	// The handler assumes nil error means "Updated Successfully".
	if err := s.Store.UpdateWorkerCPULimit(name, req.Limit); err != nil {
		if err.Error() == "Worker Not Found" {
			http.Error(w, "Worker Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Updated"))
}

// Create handles POST /workers/
func (s *Server) CreateWorkerHandler(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Limit         int    `json:"limit"`
		Name          string `json:"name"`
		ScriptContent string `json:"script_content"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}
	worker := &models.Worker{
		Name:          req.Name,
		ScriptContent: req.ScriptContent,
		CPULimit:      req.Limit,
	}

	// BUG 2 (Consumer side): The store returns nil error even if row wasn't found (0 rows affected).
	// The handler assumes nil error means "Updated Successfully".
	if err := s.Store.CreateWorker(worker); err != nil {
		http.Error(w, "Failed to create", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created"))
}
