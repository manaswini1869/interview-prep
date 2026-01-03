package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"workers/internal/model"

	"github.com/manaswini1869/debug-prepare/sc13/internal/store"
)

func MockStore() *Handler {
	return &Handler{
		Store: store.NewMemoryStore(),
	}
}

// Add Mock data to the store
func AddMockData(store *Handler) {
	store.Store.Create(model.Worker{
		ID:      "worker1",
		Script:  "script1",
		Active:  true,
		Version: 1,
	})
}

func TestDeactivateWorker(t *testing.T) {
	AddMockData(MockStore())

	// happy path
	req := httptest.NewRequest("POST", "/deactivate?id=worker1", nil)
	w := httptest.NewRecorder()
	MockStore().DeactivateWorker(w, req)
	if w.Result().StatusCode != 200 {
		t.Errorf("expected status 200, got %d", w.Result().StatusCode)
	}
	worker := MockStore().Store.Get("worker1")
	if worker.Active != false {
		t.Errorf("expected worker to be inactive")
	}

	// when id is not present
	req = httptest.NewRequest("POST", "/deactivate", nil)
	w = httptest.NewRecorder()
	MockStore().DeactivateWorker(w, req)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Result().StatusCode)
	}

	// when store returns error
	req = httptest.NewRequest("POST", "/deactivate?id=worker2", nil)
	w = httptest.NewRecorder()
	MockStore().DeactivateWorker(w, req)
	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Result().StatusCode)
	}

}
