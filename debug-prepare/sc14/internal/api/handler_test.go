package api

import (
	"net/http"
	"net/http/httptest"
	"rollout/internal/model"
	"testing"

	"github.com/manaswini1869/debug-prepare/sc14/internal/store"
)

func mockSetupHandler() *Handler {
	store := store.NewRolloutStore()
	store.Save(&model.Rollout{
		ID: "rollout1",
		Regions: map[string]string{
			"us-east-1": "completed",
			"eu-west-1": "in-progress",
		},
		Completed: false,
	})
	return &Handler{
		Store: store,
	}
}
func TestDeactivateStatus(t *testing.T) {

	handler := mockSetupHandler()
	req, err := http.NewRequest("POST", "/?id=rollout1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.DeactivateRollout(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	rollout, err := handler.Store.Get("rollout1")
	if err != nil {
		t.Fatalf("expected rollout to be found, got error: %v", err)
	}
	if !rollout.Completed {
		t.Errorf("expected rollout to be marked as completed")
	}
}
