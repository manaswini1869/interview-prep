package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type WorkerStore struct {
	Workers map[string]string
	mu      sync.RWMutex
}

// In-memory store (Simulating a database)
// INTENTIONAL BUG: Map is not thread-safe!
// var workerStore = make(map[string]string)

var workerStore *WorkerStore = &WorkerStore{
	Workers: make(map[string]string),
}

type WorkerPayload struct {
	ID     string `json:"id"`
	Script string `json:"script"`
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload WorkerPayload
	// INTENTIONAL BUG: No limit on request body size (DDoS vector)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// INTENTIONAL BUG: Race condition here under concurrent load
	workerStore.mu.Lock()
	workerStore.Workers[payload.ID] = payload.Script
	workerStore.mu.Unlock()

	// Simulate "processing" time (e.g., distributing to edge)
	time.Sleep(50 * time.Millisecond)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Deployed " + payload.ID))
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	// INTENTIONAL BUG: Concurrent read while write happens panics in Go
	workerStore.mu.RLock()
	defer workerStore.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workerStore)
}

func main() {
	http.HandleFunc("/deploy", deployHandler)
	http.HandleFunc("/list", listHandler)

	fmt.Println("Control Plane listening on :8080...")
	// INTENTIONAL BUG: ListenAndServe blocks; no graceful shutdown
	log.Fatal(http.ListenAndServe(":8080", nil))
}
