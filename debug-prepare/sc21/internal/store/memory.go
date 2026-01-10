package store

import (
	"errors"
	"sc21/internal/model"
	"sync"
)

type Store struct {
	mu     sync.RWMutex
	tokens map[string]model.Token
}

func NewStore() *Store {
	return &Store{tokens: make(map[string]model.Token)}
}

func (s *Store) Get(id string) (model.Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if token, exists := s.tokens[id]; !exists {
		return model.Token{}, errors.New("token not found")
	} else {
		return token, nil
	}
}

func (s *Store) Save(t model.Token) (model.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tokens[t.ID]; exists {
		return model.Token{}, errors.New("token already exists")
	}
	s.tokens[t.ID] = t
	return t, nil
}

func (s *Store) Revoke(id string) (model.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if token, exists := s.tokens[id]; !exists {
		return model.Token{}, errors.New("token not found")
	} else {
		token.Revoked = true
		s.tokens[id] = token
		return token, nil
	}
}
