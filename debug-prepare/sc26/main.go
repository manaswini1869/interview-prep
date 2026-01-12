package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Deployment struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"` // PENDING, RUNNING, COMPLETED, FAILED
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	data map[string]Deployment
	mu   sync.RWMutex
	// TODO: Something might be missing here for concurrent access
}

var store = Store{
	data: make(map[string]Deployment),
}

func main() {
	http.HandleFunc("/deployments", createDeployment)
	http.HandleFunc("/deployments/status", getDeployment)
	http.HandleFunc("/deployments/{id}/complete", putDeployment)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// ENDPOINT 1: Create a new deployment
// Expected behavior: Accepts JSON, stores it, returns 201.
func createDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	var d Deployment
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Simulation of validation logic
	if d.ID == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	d.CreatedAt = time.Now()
	d.Status = "PENDING"

	// BUG HINT: What happens if two requests hit this at the exact same time?
	store.data[d.ID] = d

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

// ENDPOINT 2: Get deployment status
// Expected behavior: Returns the deployment JSON if found, 404 if not.
func getDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID query parameter required", http.StatusBadRequest)
		return
	}

	// BUG HINT: Is checking for existence sufficient?
	// What does Go return when accessing a map key that doesn't exist?
	store.mu.RLock()
	defer store.mu.RUnlock()
	if _, exists := store.data[id]; !exists {
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}
	deployment := store.data[id]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployment)
	w.WriteHeader(http.StatusOK)
}

// ENDPOINT 3: Update deployment status to COMPLETED
// Expected behavior: Updates status, returns updated deployment.
func putDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID query parameter required", http.StatusBadRequest)
		return
	}

	store.mu.RLock()
	if _, exists := store.data[id]; !exists {
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}
	store.mu.RUnlock()
	store.mu.Lock()
	defer store.mu.Unlock()
	deployment := store.data[id]
	deployment.Status = "COMPLETED"
	store.data[id] = deployment

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployment)
	w.WriteHeader(http.StatusOK)
}
