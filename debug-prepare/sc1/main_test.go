package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestUpdateWorkerStatusHandler(t *testing.T) {
	// Test implementation goes here
	store = NewStore()
	testId := "test-worker-1"
	store.Save(WorkerScript{
		ID:        testId,
		Name:      "Test Worker",
		Content:   "console.log('Hello World');",
		Status:    "inactive",
		CreatedAt: 0,
	})
	payload := []byte(`{"status": "deployed"}`)
	req, _ := http.NewRequest("PUT", "/workers/"+testId+"/status", bytes.NewBuffer(payload))

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/workers/{id}/status", UpdateWorkerStatusHandler)

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	updatedWorker, _ := store.Get(testId)
	if updatedWorker.Status != "deployed" {
		t.Errorf("Worker status not updated: got %v want %v",
			updatedWorker.Status, "deployed")
	}
}
