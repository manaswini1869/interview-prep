package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	store := NewDeploymentStore()
	handler := NewHandler(store)

	r := mux.NewRouter()
	r.HandleFunc("/deployments", handler.CreateDeployment).Methods("POST")
	r.HandleFunc("/deployments", handler.ListDeployments).Methods("GET")
	r.HandleFunc("/deployments/{id}", handler.GetDeployment).Methods("GET")
	r.HandleFunc("/deployments/{id}/status", handler.UpdateStatus).Methods("PATCH")
	r.HandleFunc("/deployments/{id}", handler.DeleteDeployment).Methods("DELETE")

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
