package tests

import (
	"net/http"
	"net/http/httptest"
	"sc22/internal/handler"
	"sc22/internal/routes"
	"testing"
)

func TestEnableRoute(t *testing.T) {
	req := httptest.NewRequest("POST", "/routes/12345/enable", nil)
	w := httptest.NewRecorder()
	route := routes.Route{ID: "123", Enabled: false}
	store := routes.NewStore()
	store.Save(route)

	handler := &handler.Handler{Store: store}
	handler.EnableRoute(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204")
	}
}
