package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Route struct {
	ID         string    `json:"id"`
	Pattern    string    `json:"pattern"`
	WorkerName string    `json:"worker_name"`
	Priority   int       `json:"priority"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}

var routes = make(map[string]*Route)
var routesByWorker = make(map[string][]*Route) // worker_name -> routes
var mu sync.RWMutex

// BUG 1: This endpoint has issues with pattern validation and duplicate handling
func createRoute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pattern    string `json:"pattern"`
		WorkerName string `json:"worker_name"`
		Priority   int    `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate pattern is a valid regex
	_, err := regexp.Compile(req.Pattern)
	if err != nil {
		http.Error(w, "Invalid pattern", http.StatusBadRequest)
		return
	}

	mu.RLock()
	// Check for existing route with same pattern
	for _, route := range routes {
		if route.Pattern == req.Pattern {
			http.Error(w, "Route already exists", http.StatusConflict)
			return
		}
	}
	mu.RUnlock()

	id := fmt.Sprintf("route-%d", time.Now().Unix())
	route := &Route{
		ID:         id,
		Pattern:    req.Pattern,
		WorkerName: req.WorkerName,
		Priority:   req.Priority,
		Enabled:    true,
		CreatedAt:  time.Now(),
	}
	mu.Lock()
	routes[id] = route
	routesByWorker[req.WorkerName] = append(routesByWorker[req.WorkerName], route)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(route)
}

// BUG 2: This endpoint has issues with status code and error handling
func deleteRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	route, exists := routes[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Route Not Found"})
		return
	}
	mu.Lock()
	// Remove from routes map
	delete(routes, id)

	// Remove from worker index
	workerRoutes := routesByWorker[route.WorkerName]
	for i, r := range workerRoutes {
		if r.ID == id {
			routesByWorker[route.WorkerName] = append(workerRoutes[:i], workerRoutes[i+1:]...)
			break
		}
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Route deleted"})
}

// TODO: Implement this endpoint
// GET /api/routes/conflicts
// Should:
// 1. Find all enabled routes
// 2. Detect conflicts where multiple routes could match the same URL
//   - Two routes conflict if their patterns overlap (e.g., "/api/*" and "/api/users/*")
//   - Routes with different priorities don't conflict (higher priority wins)
//
// 3. Return groups of conflicting routes
// Response format:
//
//	{
//	  "conflicts": [
//	    {
//	      "routes": [{route1}, {route2}],
//	      "reason": "overlapping patterns with same priority"
//	    }
//	  ]
//	}

func splitPattern(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}

func normalizeSegment(seg string) string {
	if seg == "*" {
		return "*"
	}
	if strings.HasPrefix(seg, ":") {
		return ":param"
	}
	return seg
}
func patternsOverlap(p1, p2 string) bool {
	aSegs := splitPattern(p1)
	bSegs := splitPattern(p2)

	max := len(aSegs)
	if len(bSegs) > max {
		max = len(bSegs)
	}

	for i := 0; i < max; i++ {
		if i >= len(aSegs) {
			return aSegs[len(aSegs)-1] == "*"
		}
		if i >= len(bSegs) {
			return bSegs[len(bSegs)-1] == "*"
		}

		segA := normalizeSegment(aSegs[i])
		segB := normalizeSegment(bSegs[i])

		if segA == "*" || segB == "*" {
			return true
		}

		if segA != segB && segA != ":param" && segB != ":param" {
			return false
		}
	}
	return true
}

func detectConflicts(w http.ResponseWriter, r *http.Request) {
	// Implement this
	type conflicts struct {
		Routes []Route `json:"routes"`
		Reason string  `json:"reason"`
	}
	var response struct {
		Conflicts []conflicts `json:"conflicts"`
	}
	mu.RLock()
	var enabledRoutes []*Route
	for _, route := range routes {
		if route.Enabled {
			enabledRoutes = append(enabledRoutes, route)
		}
	}
	mu.RUnlock()

	for i := 0; i < len(enabledRoutes); i++ {
		for j := i + 1; j < len(enabledRoutes); j++ {
			r1 := enabledRoutes[i]
			r2 := enabledRoutes[j]

			if r1.Priority != r2.Priority {
				continue
			}

			if patternsOverlap(r1.Pattern, r2.Pattern) {
				response.Conflicts = append(response.Conflicts, conflicts{
					Routes: []Route{*r1, *r2},
					Reason: "overlapping patterns with same priority",
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/routes", createRoute).Methods("POST")
	r.HandleFunc("/api/routes/{id}", deleteRoute).Methods("DELETE")
	r.HandleFunc("/api/routes/conflicts", detectConflicts).Methods("GET")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
