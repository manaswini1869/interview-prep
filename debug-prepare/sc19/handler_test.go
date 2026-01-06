package main

import (
	"testing"
	"time"
)

func TestDeleteDeployment(t *testing.T) {
	mockDeployment := Deployment{
		ID:        "dep1",
		Script:    "echo Hello",
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	mockStore := NewDeploymentStore()
	mockStore.deployments[mockDeployment.ID] = &mockDeployment
	handler := &Handler{store: mockStore}
	expectedErr := "deployment not found"
	err := handler.store.Delete("nonexistent_id")
	if err == nil || err.Error() != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	err = handler.store.Delete(mockDeployment.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, exists := handler.store.deployments[mockDeployment.ID]; exists {
		t.Fatalf("expected deployment to be deleted")
	}
}
