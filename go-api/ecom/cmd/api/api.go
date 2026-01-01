package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/manaswini1869/interview-prep/go-api/ecom/service/user"

	"github.com/gorilla/mux"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)
	log.Println("Listening on Port", s.address)
	return http.ListenAndServe(s.address, router)
}
