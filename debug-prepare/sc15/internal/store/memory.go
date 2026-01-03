package store

import (
	"errors"
	"internal/model"
)

type Store struct {
	tokens map[string]model.Token
}

func NewStore() *Store {
	return &Store{tokens: make(map[string]model.Token)}
}

func (s *Store) Get(id string) (model.Token, error) {
	token, err := s.tokens[id]
	if !err {
		return model.Token{}, errors.New("token not found")
	}
	return token, nil
}

func (s *Store) Save(t model.Token) error {
	if _, exists := s.tokens[t.ID]; exists {
		// Update existing token
		return errors.New("token already exists")
	}
	s.tokens[t.ID] = t
	return nil
}

func (s *Store) Revoke(id string) error {
	token, exists := s.tokens[id]
	if !exists {
		return errors.New("token not found")
	}
	delete(s.tokens, token.ID)
	return nil
}
