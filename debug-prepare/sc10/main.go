package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Deployment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeploymentStore struct {
	mu          sync.RWMutex
	deployments map[string]Deployment
}

func NewDeployment() *DeploymentStore {
	return &DeploymentStore{
		deployments: make(map[string]Deployment),
	}
}

func (s *DeploymentStore) GetDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	deploymentId := r.URL.Query().Get("id")
	if deploymentId == "" {
		http.Error(w, "Missing deployment ID", http.StatusBadRequest)
		return
	}
	s.mu.RLock()
	deployment, exists := s.deployments[deploymentId]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deployment)
}

func (s *DeploymentStore) CreateDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	var deployment Deployment

	err := json.NewDecoder(r.Body).Decode(&deployment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if deployment.ID == "" || deployment.Name == "" {
		http.Error(w, "Missing deployment ID or Name", http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	if _, exists := s.deployments[deployment.ID]; exists {
		s.mu.Unlock()
		http.Error(w, "Deployment already exists", http.StatusConflict)
		return
	}
	s.deployments[deployment.ID] = deployment
	s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deployment)
}

func (s *DeploymentStore) DeleteDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	deploymentId := r.URL.Query().Get("id")
	if deploymentId == "" {
		http.Error(w, "Missing deployment ID", http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	if _, exists := s.deployments[deploymentId]; !exists {
		s.mu.Unlock()
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}
	delete(s.deployments, deploymentId)
	s.mu.Unlock()
	response := map[string]string{"message": "Deployment deleted successfully"}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(response)
}

func main() {

	store := NewDeployment()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /deployments", store.GetDeploymentHandler)
	mux.HandleFunc("POST /deployments", store.CreateDeploymentHandler)
	mux.HandleFunc("DELETE /deployments", store.DeleteDeploymentsHandler)
	http.ListenAndServe(":8080", mux)

}
