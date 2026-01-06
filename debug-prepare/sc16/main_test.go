package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestRollbackDeployment(t *testing.T) {

	deployments = make(map[string]*Deployment)
	deploymentHistory = make(map[string][]*Deployment)

	deployment1 := &Deployment{
		ID:        "deploy-1",
		AppName:   "test-app",
		Version:   "v1.0.0",
		Region:    "us-east-1",
		Status:    "completed",
		CreatedAt: time.Now().Add(-10 * time.Minute),
	}
	deployment2 := &Deployment{
		ID:        "deploy-2",
		AppName:   "test-app",
		Version:   "v1.1.0",
		Region:    "us-east-1",
		Status:    "completed",
		CreatedAt: time.Now().Add(-5 * time.Minute),
	}
	deployments[deployment1.ID] = deployment1
	deployments[deployment2.ID] = deployment2
	deploymentHistory["test-app"] = []*Deployment{deployment1, deployment2}
	expectedVersion := "v1.0.0"
	expectedStatusCode := http.StatusOK

	req, err := http.NewRequest("POST", "/api/deployments/deploy-2/rollback", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"id": "deploy-2",
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rollbackDeployment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatusCode)
	}
	if rr.Body.Len() == 0 {
		t.Errorf("handler returned empty body")
	}
	var resp Deployment
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to parse response body: %v", err)
	}
	if resp.Version != expectedVersion {
		t.Errorf("rollback did not set correct version: got %v want %v",
			resp.Version, expectedVersion)
	}
}
