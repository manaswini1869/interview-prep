package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteRoute(t *testing.T) {
	mockRoutes := []Route{
		{ID: "1", Path: "/api", Target: "service-a", Priority: 10},
		{ID: "2", Path: "/home", Target: "service-b", Priority: 5},
	}
	routes = mockRoutes

	t.Run("Route gets deleted", func(t *testing.T) {
		idToDelete := "1"
		req, err := http.NewRequest(http.MethodDelete, "/routes/"+idToDelete+"?id="+idToDelete, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		deleteRoute(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNoContent)
		}
		// Verify route is deleted
		for _, route := range routes {
			if route.ID == idToDelete {
				t.Errorf("Route with ID %s was not deleted", idToDelete)
			}
		}

	})

	t.Run("Route does not get deleted if ID not found", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodDelete, "/routes/hello?id=hello", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		deleteRoute(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
		if len(routes) != 1 {
			t.Errorf("Routes length changed unexpectedly")
		}

	})
}
