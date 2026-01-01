package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	// Assume go.mod is set up correctly
)

var store *Store

func main() {
	store = NewStore()
	r := mux.NewRouter()

	// EXISTING ENDPOINTS
	r.HandleFunc("/workers", CreateWorkerHandler).Methods("POST")
	r.HandleFunc("/workers/{id}", GetWorkerHandler).Methods("GET")
	r.HandleFunc("/workers/{id}/status", UpdateWorkerStatusHandler).Methods("PUT")

	// TODO: Your task will be to add a new endpoint here

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Control Plane Service running on :8080")
	log.Fatal(srv.ListenAndServe())
}

// --- HANDLERS ---

func CreateWorkerHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if r.Body == nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	newWorker := WorkerScript{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Content:   req.Content, // added content from request
		Status:    "provisioning",
		CreatedAt: time.Now().Unix(),
	}

	if err := store.Save(newWorker); err != nil {
		http.Error(w, "Failed to save worker", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newWorker)
}

func GetWorkerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	worker, err := store.Get(id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(worker)
}

func UpdateWorkerStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	allowed := map[string]bool{
		"provisioning": true,
		"deployed":     true,
		"failed":       true,
	}
	var req UpdateStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Status != "" && !allowed[req.Status] {
		http.Error(w, "Status parameter not supported in this endpoint", http.StatusBadRequest)
		return
	}
	err := store.Update(id, req.Status)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Status updated successfully"})
}
