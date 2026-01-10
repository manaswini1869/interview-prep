package routes

import "sync"

type Store struct {
	mu     sync.RWMutex
	routes map[string]Route
}

func NewStore() *Store {
	return &Store{routes: make(map[string]Route)}
}

func FindConflicts(routes []Route) []Route {
	// Dummy implementation for illustration
	return []Route{}
}

func (s *Store) List() []Route {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []Route
	for _, r := range s.routes {
		out = append(out, r)
	}
	return out
}

func (s *Store) Save(r Route) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[r.ID] = r
}

func (s *Store) Enable(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.routes[id]
	if !ok {
		return false
	}
	r.Enabled = true
	s.routes[id] = r
	return true
}
