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

func (s *MemoryStore) Save(route model.Route) {
	s.Routes[route.ID] = route
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
