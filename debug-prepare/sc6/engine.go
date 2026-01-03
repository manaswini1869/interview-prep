package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
)

type WorkerVersion struct {
	ID      string
	Code    string
	Version int
	Status  string
}

type VersionStore struct {
	mu       sync.Mutex
	registry map[string]WorkerVersion
}

var store = &VersionStore{
	registry: make(map[string]WorkerVersion),
}

// Intercepts the POST and PUT requests and log the HTTP Method, Path, and the User-Agent to
// stdout.
func (s *VersionStore) AuditLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the HTTP Method, Path, and User-Agent
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			log.Printf("Method: %s, Path: %s, User-Agent: %s", r.Method, r.URL.Path, r.Header.Get("User-Agent"))
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					log.Printf("Request Body: %s", string(bodyBytes))
					// Restore the io.ReadCloser to its original state
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}

			}
			next.ServeHTTP(w, r)
		}
	})
}

func (s *VersionStore) UpdateDeployment(id string, code string, version int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	curr_version := s.registry[id].Version
	if version <= curr_version {
		return errors.New("cannot deploy an older version")
	}

	// Current logic: Just overwrite whatever is there
	s.registry[id] = WorkerVersion{
		ID:      id,
		Code:    code,
		Version: version,
		Status:  "DEPLOYED",
	}
	return nil
}

func main() {

}
