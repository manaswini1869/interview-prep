package store

import (
	"errors"
	"sync"

	"github.com/manaswini1869/debug-prepare/sc20/internal/model"
)

type RolloutStore struct {
	mu       sync.Mutex
	rollouts map[string]*model.Rollout
}

func NewRolloutStore() *RolloutStore {
	return &RolloutStore{
		rollouts: make(map[string]*model.Rollout),
	}
}

func (s *RolloutStore) Get(id string) (*model.Rollout, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rollout, exists := s.rollouts[id]; !exists {
		return nil, errors.New("rollout not found")
	} else {
		return rollout, nil
	}
}

func (s *RolloutStore) Save(r *model.Rollout) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.rollouts[r.ID]; exists {
		return errors.New("duplicate ID") // or return an error indicating duplicate ID
	}
	s.rollouts[r.ID] = r
	return nil
}

func (s *RolloutStore) UpdateRegionStatus(id string, statusUpdate model.StatusUpdate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rollout, exists := s.rollouts[id]; !exists {
		return errors.New("rollout not found")
	} else {
		if rollout.Regions == nil {
			rollout.Regions = make(map[string]string)
		}
		rollout.Regions[statusUpdate.Region] = statusUpdate.Status
		return nil
	}
}
