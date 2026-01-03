package store

import (
	"errors"
	"sc14/internal/model"
	"sync"
)

type RolloutStore struct {
	mu       sync.RWMutex
	rollouts map[string]*model.Rollout
}

func NewRolloutStore() *RolloutStore {
	return &RolloutStore{
		rollouts: make(map[string]*model.Rollout),
	}
}

func (s *RolloutStore) Get(id string) (*model.Rollout, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	worker, ok := s.rollouts[id]
	if !ok {
		return nil, errors.New("rollout not found")
	}
	return worker, nil
}

func (s *RolloutStore) Save(r *model.Rollout) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rollouts[r.ID]; ok {
		return errors.New("rollout found")
	}

	s.rollouts[r.ID] = r

	return nil
}

func (s *RolloutStore) Deactivate(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	worker, ok := s.rollouts[id]
	if !ok {
		return errors.New("rollout not found")
	}

	worker.Completed = true
	s.rollouts[id] = worker

	return nil
}
