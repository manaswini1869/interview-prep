package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Route struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Target   string `json:"target"`
	Priority int    `json:"priority"`
}

var routes = []Route{
	{ID: "1", Path: "/api", Target: "service-a", Priority: 10},
	{ID: "2", Path: "/home", Target: "service-b", Priority: 5},
}

func main() {
	http.HandleFunc("/routes", handleRoutes)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRoutes(w, r)
	case http.MethodPost:
		createRoute(w, r)
	case http.MethodDelete:
		deleteRoute(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ENDPOINT 1: List routes
// Expected behavior: Return all routes.
// BUG: The user wants to filter by Target, e.g. /routes?target=service-a
// Currently, it ignores the query param.
func getRoutes(w http.ResponseWriter, r *http.Request) {
	targetFilter := r.URL.Query().Get("target")
	if targetFilter == "" {
		http.Error(w, "Bad request: missing target filter", http.StatusBadRequest)
		return
	}
	var result []Route
	// Logic is missing here to actually filter
	for _, route := range routes {
		if route.Target == targetFilter {
			result = append(result, route)
		}
	}
	if result == nil {
		result = []Route{}
	}

	// Challenge: If targetFilter is present but no routes match, return empty list [] not null
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)
}

// ENDPOINT 2: Create Route
// Expected behavior: Add a route. prevent duplicate Paths.
func createRoute(w http.ResponseWriter, r *http.Request) {
	var newRoute Route
	if err := json.NewDecoder(r.Body).Decode(&newRoute); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// BUG: The loop below logic is slightly flawed for updating/checking duplicates
	for _, route := range routes {
		if route.Path == newRoute.Path || route.ID == newRoute.ID || route.Target == newRoute.Target {
			http.Error(w, "Path already exists", http.StatusConflict)
			// Missing return statement? What happens here?
			return
		}
	}

	routes = append(routes, newRoute)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newRoute)
	w.WriteHeader(http.StatusCreated)
}

// ENDPOINT 3: Delete Route
// Expected behavior: Delete a route by ID.
func deleteRoute(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID query parameter required", http.StatusBadRequest)
		return
	}

	if routes == nil {
		http.Error(w, "No routes available", http.StatusNotFound)
		return
	}

	for i, route := range routes {
		if route.ID == id {
			routes = append(routes[:i], routes[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Route not found", http.StatusNotFound)
}
