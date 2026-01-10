package api

import (
	"net/http"
	"net/http/httptest"
	"sc21/internal/model"
	"sc21/internal/store"
	"strings"
	"testing"
)

func TestRevokeToken(t *testing.T) {
	store := store.NewStore()
	handler := &Handler{Store: store}
	// Prepopulate the store with a token
	token := model.Token{
		ID:      "token123",
		Revoked: false,
		Scopes:  []string{"read", "write"},
	}
	_, err := store.Save(token)
	if err != nil {
		t.Fatalf("Failed to save token: %v", err)
	}
	t.Run("Revoke existing token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/revoke?id=token123", nil)
		w := httptest.NewRecorder()
		handler.RevokeToken(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status 204, got %d", resp.StatusCode)
		}
		if token, _ := store.Get("token123"); !token.Revoked {
			t.Fatalf("Expected token to be revoked")
		}
	})
	t.Run("Revoke non-existing token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/revoke?id=token_no", nil)
		w := httptest.NewRecorder()
		handler.RevokeToken(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status 404, got %d", resp.StatusCode)
		}
		if !strings.Contains(w.Body.String(), "token not found") {
			t.Fatalf("Unexpected response body: %s", w.Body.String())
		}

	})

}
