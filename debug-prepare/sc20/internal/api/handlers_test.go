package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/manaswini1869/debug-prepare/sc20/internal/model"
	"github.com/manaswini1869/debug-prepare/sc20/internal/store"
)

func TestUpdateRolloutRegionStatus(t *testing.T) {
	mockStore := store.NewRolloutStore()
	handler := &Handler{Store: mockStore}

	exisintgRollout := &model.Rollout{
		ID:      "rollout1",
		Regions: map[string]string{"us-east-1": "pending", "eu-west-1": "pending"},
	}
	err := mockStore.Save(exisintgRollout)
	if err != nil {
		t.Fatalf("failed to setup mock store: %v", err)
	}
	t.Run("Update the existing regions", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/rollout/update-region-status?id=rollout1", strings.NewReader(`{"region":"us-east-1","status":"completed"}`))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateRegionStatus(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}
		updatedRollout, _ := mockStore.Get("rollout1")
		if updatedRollout.Regions["us-east-1"] != "completed" {
			t.Errorf("region status not updated: got %v want %v", updatedRollout.Regions["us-east-1"], "completed")
		}
	})

	t.Run("Update the non existing region", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/rollout/update-region-status?id=rollout_no", strings.NewReader(`{"region":"us-east-1","status":"completed"}`))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateRegionStatus(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
		if !strings.Contains(rr.Body.String(), "rollout not found") {
			t.Errorf("unexpected response body: got %v", rr.Body.String())
		}

	})

}
