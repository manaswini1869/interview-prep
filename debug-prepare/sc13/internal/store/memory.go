package store

import (
	"errors"
	"sync"
	"workers/internal/model"
)

type MemoryStore struct {
	mu      sync.RWMutex
	workers map[string]model.Worker
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		workers: make(map[string]model.Worker),
	}
}

func (s *MemoryStore) Create(w model.Worker) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.workers[w.ID]; exists {
		return errors.New("worker already exists")
	}

	s.workers[w.ID] = w
	return nil
}

func (s *MemoryStore) Get(id string) (model.Worker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w, ok := s.workers[id]
	if !ok {
		return model.Worker{}, errors.New("not found")
	}
	return w, nil
}

func (s *MemoryStore) DeactivateWorker(id string) (model.Worker, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.workers[id]
	if !ok {
		return model.Worker{}, errors.New("not found")
	}
	w.Active = false
	s.workers[id] = w
	return w, nil
}

func (s *MemoryStore) List() []model.Worker {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := []model.Worker{}
	for _, w := range s.workers {
		result = append(result, w)
	}
	return result
}
