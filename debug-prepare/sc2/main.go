package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	shareCounter := NewContainer()
	Server := &Server{container: shareCounter}

	r.HandleFunc("/stats/{id}/record", Server.RecordHitHandler).Methods("POST")
	r.Handle("/stats/{id}", AuthMiddleware(http.HandlerFunc(Server.GetStatsHandler))).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8081",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	srv.ListenAndServe()
}

type Server struct {
	container *Container
}

func NewContainer() *Container {
	return &Container{
		statsMap: make(map[string]int),
	}
}

func (s *Server) RecordHitHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Simulate some processing latency
	time.Sleep(50 * time.Millisecond)

	s.container.RecordHit(id)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetStatsHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	count := s.container.GetStats(id)
	resp := map[string]int{"count": count}
	json.NewEncoder(w).Encode(resp)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token != "admin-secret" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
