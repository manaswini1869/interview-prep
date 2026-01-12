package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPutDeployment(t *testing.T) {
	// test data
	dep1 := Deployment{
		ID:        "dep1",
		Status:    "RUNNING",
		CreatedAt: time.Now(),
	}
	mockStore := Store{
		data: map[string]Deployment{
			"dep1": dep1,
		},
	}
	t.Run("Deployment Updated", func(t *testing.T) {
		store = mockStore // inject mock store
		req, err := http.NewRequest("PUT", "/deployments/dep1/complete?id=dep1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(putDeployment)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		expectedStatus := "COMPLETED"
		updatedDep, _ := store.data["dep1"]
		if updatedDep.Status != expectedStatus {
			t.Errorf("Deployment status not updated: got %v want %v",
				updatedDep.Status, expectedStatus)
		}
	})
	t.Run("Deployment Updated", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/deployments/dep/complete?id=dep", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(putDeployment)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})

}
