package main

import (
	"fmt"
	"sync"
	"time"
)

type Deployment struct {
	ID     string
	Region string
	Status string
}

type DeploymentStore struct {
	mu          sync.Mutex
	deployments map[string]Deployment
}

var store = &DeploymentStore{
	deployments: make(map[string]Deployment),
}

// Simulates a call to an external region validation service
func validateRegionStatus(region string) bool {
	time.Sleep(200 * time.Millisecond) // This is slow!
	return true
}

func (s *DeploymentStore) Add(d Deployment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deployments[d.ID] = d
}

func (s *DeploymentStore) ListActive() []Deployment {
	s.mu.Lock()
	tempList := make([]Deployment, 0, len(s.deployments))
	for _, d := range s.deployments {
		tempList = append(tempList, d)
	}
	s.mu.Unlock()

	results := []Deployment{}
	for _, d := range tempList {
		// We only want to return deployments in healthy regions
		if validateRegionStatus(d.Region) {
			results = append(results, d)
		}
	}
	return results
}

func (s *DeploymentStore) DeleteDeployment(id string) error {
	if _, exists := s.deployments[id]; !exists {
		return fmt.Errorf("deployment not found")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.deployments, id)
	return nil
}
