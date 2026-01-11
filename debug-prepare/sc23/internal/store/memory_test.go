package store

import (
	"sc23/internal/model"
	"testing"
)

func TestGetConflicts(t *testing.T) {
	mockStore := NewMemoryStore()
	route1 := model.Route{
		ID:       "1",
		Pattern:  "/api/v1/resource",
		Priority: 1,
		Enabled:  true,
	}
	route2 := model.Route{
		ID:       "2",
		Pattern:  "/api/v1/resource",
		Priority: 1,
		Enabled:  true,
	}
	route3 := model.Route{
		ID:       "3",
		Pattern:  "/api/v1/other",
		Priority: 2,
		Enabled:  true,
	}
	mockStore.Save(route1)
	mockStore.Save(route2)
	mockStore.Save(route3)

	t.Run("Get Conflicts with conflicts", func(t *testing.T) {
		mockResult := []model.Route{route1, route2}
		conflicts := mockStore.GetConflicts()
		if len(conflicts) != len(mockResult) {
			t.Errorf("Expected %d conflicts, got %d", len(mockResult), len(conflicts))
		}
		for i, r := range conflicts {
			if r != mockResult[i] {
				t.Errorf("Expected route %v, got %v", mockResult[i], r)
			}
		}
	})

	t.Run("Get Conflicts with no conflicts", func(t *testing.T) {
		mockResult := []model.Route{}

		delete(mockStore.Routes, "2")
		conflicts := mockStore.GetConflicts()
		if len(conflicts) != len(mockResult) {
			t.Errorf("Expected %d conflicts, got %d", len(mockResult), len(conflicts))
		}
		for i, r := range conflicts {
			if r != mockResult[i] {
				t.Errorf("Expected route %v, got %v", mockResult[i], r)
			}
		}
	})
}
