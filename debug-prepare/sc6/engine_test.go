package main

import (
	"sync"
	"testing"
)

func setup() {
	store = &VersionStore{
		registry: make(map[string]WorkerVersion),
	}
	store.registry["worker1"] = WorkerVersion{
		ID:      "worker1",
		Code:    "initial code",
		Version: 10,
		Status:  "DEPLOYED",
	}
}

func OlderVersionError() error {
	return store.UpdateDeployment("worker1", "hello", 10)
}

func NewerVersionError() error {
	return store.UpdateDeployment("worker1", "helloworld", 11)
}

func TestUpdateDeployment(t *testing.T) {
	setup()

	var wg sync.WaitGroup
	wg.Add(2)
	var err1, err2 error

	go func() {
		defer wg.Done()
		err1 = OlderVersionError()
	}()
	go func() {
		defer wg.Done()
		err2 = NewerVersionError()
	}()
	wg.Wait()

	if err1 == nil && err2 == nil {
		t.Errorf("Expected one of the updates to fail due to version conflict, but both succeeded")
	}

	// Allow some time for goroutines to finish
	// In real tests, use sync.WaitGroup or channels for synchronization
	// Here we just use a simple sleep for demonstration purposes
	// time.Sleep(1 * time.Second)
	finalVersion := store.registry["worker1"].Version
	if finalVersion != 11 {
		t.Errorf("Expected final version to be 11, got %d", finalVersion)
	}
}
