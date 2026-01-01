package main

import (
	"database/sql"
	"log"

	"github.com/manaswini1869/interview-prep/go-api/ecom/cmd/api"
	"github.com/manaswini1869/interview-prep/go-api/ecom/db"
)

func main() {
	db, err := db.NewMySQLStorage()

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

func initStorage(db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully Connected")
}
