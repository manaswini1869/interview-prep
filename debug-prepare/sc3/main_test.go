package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setup() {
	// Returns a fresh store for every single test
	store = &DeploymentStore{
		deployments: make(map[string]Deployment),
	}
	store.Add(Deployment{ID: "1", Region: "us-east", Status: "active"})
	store.Add(Deployment{ID: "2", Region: "eu-west", Status: "active"})
}

func TestHappyPathDeleteDeployment(t *testing.T) {
	// Start the server
	setup()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/deployments/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)
	DeleteDeploymentHandler(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	store.mu.Lock()
	_, exists := store.deployments["1"]
	store.mu.Unlock()

	if exists {
		t.Errorf("expected deployment 1 to be deleted from global store")
	}

}

func TestDeleteDeploymentInvalidID(t *testing.T) {
	// Start the server in a goroutine
	setup()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/deployments/6", nil)
	vars := map[string]string{
		"id": "6",
	}
	req = mux.SetURLVars(req, vars)
	DeleteDeploymentHandler(rr, req)

	if status := rr.Code; status != 404 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 404)
	}
}
