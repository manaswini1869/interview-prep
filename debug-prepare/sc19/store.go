package main

import (
	"errors"
	"sync"
)

type DeploymentStore struct {
	mu          sync.RWMutex
	deployments map[string]*Deployment
}

func NewDeploymentStore() *DeploymentStore {
	return &DeploymentStore{
		deployments: make(map[string]*Deployment),
	}
}

func (s *DeploymentStore) Create(d *Deployment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deployments[d.ID] = d
}

func (s *DeploymentStore) Get(id string) (*Deployment, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	d, ok := s.deployments[id]
	return d, ok
}

func (s *DeploymentStore) List() []*Deployment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.deployments) == 0 {
		return []*Deployment{}
	}

	out := []*Deployment{}
	for _, d := range s.deployments {
		out = append(out, d)
	}
	return out
}

func (s *DeploymentStore) UpdateStatus(id string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d, ok := s.deployments[id]; ok {
		d.Status = status
		return nil
	} else {
		// Handle the case where the deployment does not exist
		return errors.New("deployment not found")
	}
}

func (s *DeploymentStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d, ok := s.deployments[id]; ok {
		delete(s.deployments, d.ID)
		return nil
	} else {
		// Handle the case where the deployment does not exist
		return errors.New("deployment not found")
	}
}
