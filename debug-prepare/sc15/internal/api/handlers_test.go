package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manaswini1869/debug-prepare/sc15/internal/store"
)

func TestRevokeToken(t *testing.T) {

	resp, err := http.NewRequest("POST", "/?id=123", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := &Handler{Store: store.NewStore()}
	handler.RevokeToken(rr, resp)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	if handler.Store.tokens["123"] != nil {
		t.Errorf("handler did not revoke the token")
	}

}
