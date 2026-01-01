package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewMySQLStorage() (*sql.DB, error) {
	host := "localhost"
	port := 5432
	user := "root"
	password := "something"
	dbname := "ecom"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
