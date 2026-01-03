package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/mux"
)

func TestUpdateWorkerStatusHandler(t *testing.T) {
	server := &Server{
		container: NewContainer(),
	}
	router := mux.NewRouter()
	router.HandleFunc("/stats/{id}/record", server.RecordHitHandler).Methods("POST")
	router.HandleFunc("/stats/{id}", server.GetStatsHandler).Methods("GET")

	workerID := "worker1"
	const concurrentHits = 100
	var wg sync.WaitGroup

	// 1. Fire off concurrent POST requests
	wg.Add(concurrentHits)
	for i := 0; i < concurrentHits; i++ {
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("POST", "/stats/"+workerID+"/record", nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("POST returned %v", rr.Code)
			}
		}()
	}

	// 2. Wait for all goroutines to finish
	wg.Wait()

	// 3. Test getting stats (Verification)
	req, _ := http.NewRequest("GET", "/stats/"+workerID, nil)
	req.Header.Set("X-Auth-Token", "admin-token")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET returned %v", status)
	}

	var resp map[string]int
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["count"] != concurrentHits {
		t.Errorf("Race condition or logic error: got %v want %v", resp["count"], concurrentHits)
	}
}

func TestMiddlewareAuth(t *testing.T) {
	server := &Server{container: NewContainer()}
	router := mux.NewRouter()
	router.HandleFunc("/stats/{id}", server.GetStatsHandler).Methods("GET")
	req, _ := http.NewRequest("GET", "/stats/worker1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
	// Now test with the correct header
	req, _ = http.NewRequest("GET", "/stats/worker1", nil)
	req.Header.Set("X-Auth-Token", "admin-token")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println("Middleware auth tests passed")
}
