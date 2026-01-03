package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var store = NewDeployment()

func storeInitialization() {
	store.deployments["dep1"] = Deployment{ID: "dep1", Name: "Deployment One"}
	store.deployments["dep2"] = Deployment{ID: "dep2", Name: "Deployment Two"}

}

func TestDeleteDeploymentHandler(t *testing.T) {
	storeInitialization()
	// happy path
	req, _ := http.NewRequest("DELETE", "/deployments?id=dep1", nil)
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/deployments", store.DeleteDeploymentsHandler)
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
	// missing deployment ID
	req, _ = http.NewRequest("DELETE", "/deployments", nil)
	rr = httptest.NewRecorder()
	r = mux.NewRouter()
	r.HandleFunc("/deployments", store.DeleteDeploymentsHandler)
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for missing ID: got %v want %v",
			status, http.StatusBadRequest)
	}

	// deployment not found
	req, _ = http.NewRequest("DELETE", "/deployments?id=nonexistent", nil)
	rr = httptest.NewRecorder()
	r = mux.NewRouter()
	r.HandleFunc("/deployments", store.DeleteDeploymentsHandler)
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for not found: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestGetDeploymentHandler(t *testing.T) {
	storeInitialization()
	// happy path
	req, _ := http.NewRequest("GET", "/deployments?id=dep1", nil)
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/deployments", store.GetDeploymentHandler)
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// missing deployment ID
	req, _ = http.NewRequest("GET", "/deployments", nil)
	rr = httptest.NewRecorder()
	r = mux.NewRouter()
	r.HandleFunc("/deployments", store.GetDeploymentHandler)
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for missing ID: got %v want %v",
			status, http.StatusBadRequest)
	}

	// deployment not found
	req, _ = http.NewRequest("GET", "/deployments? id=nonexistent", nil)
	rr = httptest.NewRecorder()
	r = mux.NewRouter()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for not found: got %v want %v",
			status, http.StatusNotFound)
	}

}

func TestCreateDeploymentHandler(t *testing.T) {
	storeInitialization()
	// happy path
	reqBody := `{"id": "dep3", "name": "Deployment Three"}`
	req, _ := http.NewRequest("POST", "/deployments", bytes.NewBuffer([]byte(reqBody)))
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/deployments", store.CreateDeploymentHandler)
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// missing deployment ID or Name
	reqBody = `{"id": "", "name": "Deployment Four"}`
	req, _ = http.NewRequest("POST", "/deployments", bytes.NewBuffer([]byte(reqBody)))
	rr = httptest.NewRecorder()
	r = mux.NewRouter()
	r.HandleFunc("/deployments", store.CreateDeploymentHandler)
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for missing fields: got %v want %v",
			status, http.StatusBadRequest)
	}

	// deployment already exists
	reqBody = `{"id": "dep1", "name": "Deployment One"}`
	req, _ = http.NewRequest("POST", "/deployments", bytes.NewBuffer([]byte(reqBody)))
	rr = httptest.NewRecorder()
	r = mux.NewRouter()
	r.HandleFunc("/deployments", store.CreateDeploymentHandler)
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code for existing deployment: got %v want %v",
			status, http.StatusConflict)
	}
}
