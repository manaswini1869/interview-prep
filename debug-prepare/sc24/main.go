package main

import (
	"log"
	"net/http"
	"sc24/api"
	"sc24/store"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Connection string for Docker Compose
	connStr := "user=user password=password dbname=worker_db sslmode=disable host=localhost"

	// Retry logic for DB connection (simulation helper)
	var dbStore *store.Store
	var err error
	for i := 0; i < 5; i++ {
		dbStore, err = store.NewStore(connStr)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
		log.Println("Waiting for DB...")
	}
	if err != nil {
		log.Fatal(err)
	}

	server := &api.Server{Store: dbStore}
	r := mux.NewRouter()

	r.HandleFunc("/workers/{name}", server.GetWorkerHandler).Methods("GET")
	r.HandleFunc("/workers/{name}/cpu", server.UpdateCPUHandler).Methods("PUT")
	r.HandleFunc("/workers", server.CreateWorkerHandler).Methods("PUT")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
