package store

import "sc23/internal/model"

type MemoryStore struct {
	Routes map[string]model.Route
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Routes: make(map[string]model.Route),
	}
}

func (s *MemoryStore) Save(route model.Route) (model.Route, bool) {
	if _, exists := s.Routes[route.ID]; exists {
		return model.Route{}, false
	}

	s.Routes[route.ID] = route
	return s.Routes[route.ID], true
}

func (s *MemoryStore) Get(id string) (model.Route, bool) {
	route, ok := s.Routes[id]
	return route, ok
}

func (s *MemoryStore) List() []model.Route {
	var res []model.Route
	for _, r := range s.Routes {
		res = append(res, r)
	}
	return res
}

func (s *MemoryStore) GetConflicts() []model.Route {
	routes := s.List()
	conflicts := [][]model.Route{}

	// get conflicts only for enabled routes based on pattern and priority
	for i := 0; i < len(routes); i++ {
		for j := i + 1; j < len(routes); j++ {
			if routes[i].Enabled && routes[j].Enabled &&
				routes[i].Pattern == routes[j].Pattern &&
				routes[i].Priority == routes[j].Priority {
				conflicts = append(conflicts, []model.Route{routes[i], routes[j]})
			}
		}
	}
	var flatConflicts []model.Route
	for _, pair := range conflicts {
		flatConflicts = append(flatConflicts, pair...)
	}
	return flatConflicts

}
