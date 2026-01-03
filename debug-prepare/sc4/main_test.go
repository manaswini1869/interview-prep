package main

import (
	"context"
	"testing"
)

func TestPostRollbackDeployment(t *testing.T) {
	store := &MockStore{
		RollbackDeployment: func(ctx context.Context, serviceID, deploymentID string) error {
			return nil
		},
	}

	handler := &Handler{store: store}
	req := httptest.NewRequest("POST", "/services/svc1/deployments/deploy1/rollback", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("service_id", "svc1")
	rctx.URLParams.Add("deployment_id", "deploy1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()
	handler.PostRollbackDeployment(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("expected status %v, got %v", http.StatusNoContent, status)
	}

}
