package main

import (
	"errors"
	"sync"
)

// In-memory simulation of a database
type Store struct {
	mu   sync.RWMutex
	data map[string]WorkerScript
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]WorkerScript),
	}
}

func (s *Store) Save(w WorkerScript) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[w.ID]; exists {
		return errors.New("worker with this ID already exists")
	}
	s.data[w.ID] = w
	return nil
}

func (s *Store) Get(id string) (*WorkerScript, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[id]
	if !ok {
		// Simulating a database "not found" behavior
		return nil, nil
	}
	return &val, nil
}

func (s *Store) Update(id string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[id]
	val.Status = status
	s.data[id] = val
	if !ok {
		// Simulating a database "not found" behavior
		return nil
	}
	return nil
}
