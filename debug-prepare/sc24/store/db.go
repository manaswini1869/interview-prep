package store

import (
	"context"
	"database/sql"
	"fmt"
	"sc24/models"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewTestStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func NewStore(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// GetWorkerByName retrieves a worker from the DB
func (s *Store) GetWorkerByName(name string) (*models.Worker, error) {
	w := &models.Worker{}
	// INTENTIONAL BUG 1: The query is correct, but the error handling ignores ErrNoRows specifically
	// or the Scan mapping might be slightly off in a real scenario.
	// For this simulation: We are scanning into the wrong variable or missing a check.
	err := s.db.QueryRowContext(context.Background(), "SELECT id, name, script_content, cpu_limit FROM workers WHERE name = $1", name).Scan(&w.ID, &w.Name, &w.ScriptContent, &w.CPULimit)

	if err != nil {
		// ISSUE: If err is sql.ErrNoRows, we return it directly.
		// The handler might not be checking this specific error to return 404.
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sql: no rows in result set")
		}
		return nil, err
	}
	return w, nil
}

// UpdateWorkerCPULimit updates the CPU limit for a worker
func (s *Store) UpdateWorkerCPULimit(name string, newLimit int) error {
	// INTENTIONAL BUG 2: Logic error / SQL Syntax error
	// We are trying to update, but if the worker doesn't exist, this Exec won't return an error,
	// it will just return 0 rows affected. The system expects an error if the worker is missing.
	res, err := s.db.Exec("UPDATE workers SET cpu_limit = $1 WHERE name = $2", newLimit, name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return fmt.Errorf("Worker Not Found")
		}
		return err
	}

	// Hint: We should check res.RowsAffected() here, but we aren't.
	// This implies a successful 200 OK update even if the worker name was a typo.
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("Worker Not Found")
	}
	return nil
}

// CreateWorker inserts a new worker
// This helper is fine, provided for the test context.
func (s *Store) CreateWorker(w *models.Worker) error {
	_, err := s.db.Exec("INSERT INTO workers (name, script_content, cpu_limit) VALUES ($1, $2, $3)",
		w.Name, w.ScriptContent, w.CPULimit)

	return err
}
