package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Deployment struct {
	ID          string     `json:"id"`
	AppName     string     `json:"app_name"`
	Version     string     `json:"version"`
	Region      string     `json:"region"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

var deployments = make(map[string]*Deployment)
var deploymentHistory = make(map[string][]*Deployment) // app_name -> deployments

// BUG 1: This endpoint has issues with status updates and response handling
func createDeployment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AppName string `json:"app_name"`
		Version string `json:"version"`
		Region  string `json:"region"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := fmt.Sprintf("deploy-%d", time.Now().Unix())
	if _, exists := deployments[id]; exists {
		http.Error(w, "Deployment ID already exists", http.StatusConflict)
		return
	}
	deployment := &Deployment{
		ID:        id,
		AppName:   req.AppName,
		Version:   req.Version,
		Region:    req.Region,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	deployments[id] = deployment
	deploymentHistory[req.AppName] = append(deploymentHistory[req.AppName], deployment)

	// Simulate async deployment
	go func() {
		time.Sleep(2 * time.Second)
		deployment.Status = "completed"
		now := time.Now()
		deployment.CompletedAt = &now
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deployment)
}

// BUG 2: This endpoint has issues with filtering and handling missing deployments
func getDeploymentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing deployment ID", http.StatusBadRequest)
		return
	}
	deployment, exists := deployments[id]
	if !exists {
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deployment)
}

// TODO: Implement this endpoint
// POST /api/deployments/{id}/rollback
// Should:
// 1. Find the deployment by ID
// 2. Find the previous successful deployment for the same app
// 3. Create a new deployment with the previous version
// 4. Return the new deployment
// Handle errors: deployment not found, no previous version, app not found
func rollbackDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing deployment ID", http.StatusBadRequest)
		return
	}
	current_deploy, exists := deployments[id]
	if !exists {
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}

	history, exists := deploymentHistory[current_deploy.AppName]
	if !exists || len(history) < 2 {
		http.Error(w, "No previous deployment to rollback to", http.StatusNotFound)
		return
	}
	var previous_deploy *Deployment
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].ID == id {
			if i > 0 {
				previous_deploy = history[i-1]
			}
			break
		}
	}
	if previous_deploy == nil || previous_deploy.Status != "completed" {
		http.Error(w, "No previous successful deployment to rollback to", http.StatusNotFound)
		return
	}
	new_id := fmt.Sprintf("deploy-%d", time.Now().Unix())
	new_deployment := &Deployment{
		ID:        new_id,
		AppName:   current_deploy.AppName,
		Version:   previous_deploy.Version,
		Region:    current_deploy.Region,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	deployments[new_id] = new_deployment
	deploymentHistory[current_deploy.AppName] = append(deploymentHistory[current_deploy.AppName], new_deployment)

	go func() {
		time.Sleep(2 * time.Second)
		new_deployment.Status = "completed"
		now := time.Now()
		new_deployment.CompletedAt = &now
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(new_deployment)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	// Implement this
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/deployments", createDeployment).Methods("POST")
	r.HandleFunc("/api/deployments/{id}", getDeploymentStatus).Methods("GET")
	r.HandleFunc("/api/deployments/{id}/rollback", rollbackDeployment).Methods("POST")
	r.HandleFunc("/api/health", healthCheck).Methods("GET")
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
