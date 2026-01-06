// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Key       string            `json:"key"`
	Value     string            `json:"value"`
	Version   int               `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	UpdatedAt time.Time         `json:"updated_at"`
}

var configs = make(map[string]*Config) // namespace:key -> Config
var mu sync.Mutex

// BUG 1: This endpoint has concurrency issues and version tracking problems
func updateConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	key := vars["key"]

	if key == "" || namespace == "" {
		http.Error(w, "Namespace and key are required", http.StatusBadRequest)
		return
	}

	var req struct {
		Value    string            `json:"value"`
		Metadata map[string]string `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	configKey := fmt.Sprintf("%s:%s", namespace, key)
	existing := configs[configKey]

	var version int = 1
	if existing != nil {
		version = existing.Version
	}

	config := &Config{
		ID:        configKey,
		Namespace: namespace,
		Key:       key,
		Value:     req.Value,
		Version:   version + 1,
		Metadata:  req.Metadata,
		UpdatedAt: time.Now(),
	}

	configs[configKey] = config

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(config)
}

// BUG 2: This endpoint has issues with namespace filtering and response format
func listConfigs(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")

	if namespace == "" {
		http.Error(w, "Namespace query parameter is required", http.StatusBadRequest)
		return
	}
	var result []*Config
	for key, config := range configs {
		if namespace != "" {
			result = append(result, config)
		} else if strings.HasPrefix(key, namespace) {
			result = append(result, config)
		}
	}
	if len(result) == 0 {
		http.Error(w, "No configs found in the specified namespace", http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// TODO: Implement this endpoint
// GET /api/configs/compare?namespace1={ns1}&namespace2={ns2}
// Should:
// 1. Get all configs from namespace1 and namespace2
// 2. Compare them and return:
//   - configs only in namespace1
//   - configs only in namespace2
//   - configs in both with different values
//   - configs in both with same values
//
// Response format:
//
//	{
//	  "only_in_ns1": [...],
//	  "only_in_ns2": [...],
//	  "different_values": [{...}],
//	  "same_values": [...]
//	}

type CompareResult struct {
	OnlyInNS1       []*Config    `json:"only_in_ns1"`
	OnlyInNS2       []*Config    `json:"only_in_ns2"`
	DifferentValues [][2]*Config `json:"different_values"`
	SameValues      []*Config    `json:"same_values"`
}

func compareConfigs(w http.ResponseWriter, r *http.Request) {
	// Implement this
	namespace1 := r.URL.Query().Get("namespace1")
	namespace2 := r.URL.Query().Get("namespace2")

	if namespace1 == "" || namespace2 == "" {
		http.Error(w, "Namespace query parameter is required", http.StatusBadRequest)
		return
	}
	var InNS1 []*Config
	var InNS2 []*Config
	mu.Lock()
	for key, config := range configs {
		if strings.HasPrefix(key, namespace1) {
			InNS1 = append(InNS1, config)
		} else if strings.HasPrefix(key, namespace2) {
			InNS2 = append(InNS2, config)
		}
	}
	mu.Unlock()
	if len(InNS1) == 0 && len(InNS2) == 0 {
		http.Error(w, "No configs found in the specified namespaces", http.StatusNotFound)
		return
	}
	ns2ByKey := make(map[string]*Config)
	for _, c := range InNS2 {
		ns2ByKey[c.Key] = c
	}
	var (
		OnlyInNS1       []*Config
		OnlyInNS2       []*Config
		DifferentValues [][2]*Config
		SameValues      []*Config
	)

	// Compare namespace1 against namespace2
	for _, c1 := range InNS1 {
		if c2, ok := ns2ByKey[c1.Key]; ok {
			if c1.Value == c2.Value && c1.Version == c2.Version {
				SameValues = append(SameValues, c1)
			} else {
				DifferentValues = append(DifferentValues, [2]*Config{c1, c2})
			}
			delete(ns2ByKey, c1.Key)
		} else {
			OnlyInNS1 = append(OnlyInNS1, c1)
		}
	}

	for _, c2 := range ns2ByKey {
		OnlyInNS2 = append(OnlyInNS2, c2)
	}

	result := CompareResult{
		OnlyInNS1:       OnlyInNS1,
		OnlyInNS2:       OnlyInNS2,
		DifferentValues: DifferentValues,
		SameValues:      SameValues,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/configs/{namespace}/{key}", updateConfig).Methods("PUT")
	r.HandleFunc("/api/configs/", listConfigs).Methods("GET")
	r.HandleFunc("/api/configs/compare", compareConfigs).Methods("GET")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
