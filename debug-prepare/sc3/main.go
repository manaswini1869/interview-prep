package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Seed some data
	store.Add(Deployment{ID: "1", Region: "us-east", Status: "active"})
	store.Add(Deployment{ID: "2", Region: "eu-west", Status: "active"})
	store.Add(Deployment{ID: "3", Region: "ap-south", Status: "active"})

	r.HandleFunc("/deployments", ListDeploymentsHandler).Methods("GET")
	// Assume there is a POST handler elsewhere...
	r.HandleFunc("/deployments/{id}", DeleteDeploymentHandler).Methods("DELETE")

	http.ListenAndServe(":8082", r)
}

func ListDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	deps := store.ListActive()
	json.NewEncoder(w).Encode(deps)
}

func DeleteDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid Deployment ID", http.StatusBadRequest)
		return
	}

	err := store.DeleteDeployment(id)
	if err != nil {
		http.Error(w, "Deployment Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
