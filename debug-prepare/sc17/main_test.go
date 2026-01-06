package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompareConfigs(t *testing.T) {
	// happy path test
	configs["ns1:key1"] = &Config{Key: "key1", Namespace: "ns1", Value: "value1", Version: 1}
	configs["ns1:key2"] = &Config{Key: "key2", Namespace: "ns1", Value: "value2", Version: 1}
	configs["ns1:key3"] = &Config{Key: "key3", Namespace: "ns1", Value: "value3", Version: 1}
	configs["ns2:key1"] = &Config{Key: "key1", Namespace: "ns2", Value: "value2", Version: 1}
	configs["ns2:key2"] = &Config{Key: "key2", Namespace: "ns2", Value: "value2", Version: 1}
	configs["ns2:key5"] = &Config{Key: "key5", Namespace: "ns2", Value: "value5", Version: 1}

	expectedResponse := &CompareResult{
		OnlyInNS1: []*Config{
			{Key: "key3", Namespace: "ns1", Value: "value3", Version: 1},
		},
		OnlyInNS2: []*Config{
			{Key: "key5", Namespace: "ns2", Value: "value5", Version: 1},
		},
		DifferentValues: [][2]*Config{
			{
				{Key: "key1", Namespace: "ns1", Value: "value1", Version: 1},
				{Key: "key1", Namespace: "ns2", Value: "value2", Version: 1},
			},
		},
		SameValues: []*Config{
			{Key: "key2", Namespace: "ns1", Value: "value2", Version: 1},
		},
	}
	req, err := http.NewRequest("GET", "/api/configs/compare?namespace1=ns1&namespace2=ns2", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(compareConfigs)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var resp CompareResult
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to parse response body: %v", err)
	}
	if len(resp.OnlyInNS1) != len(expectedResponse.OnlyInNS1) || resp.OnlyInNS1[0].Key != expectedResponse.OnlyInNS1[0].Key {
		t.Errorf("OnlyInNS1 length mismatch: got %v want %v", len(resp.OnlyInNS1), len(expectedResponse.OnlyInNS1))
	}
	if len(resp.OnlyInNS2) != len(expectedResponse.OnlyInNS2) || resp.OnlyInNS2[0].Key != expectedResponse.OnlyInNS2[0].Key {
		t.Errorf("OnlyInNS2 length mismatch: got %v want %v", len(resp.OnlyInNS2), len(expectedResponse.OnlyInNS2))
	}
	if len(resp.DifferentValues) != len(expectedResponse.DifferentValues) || resp.DifferentValues[0][0].Key != expectedResponse.DifferentValues[0][0].Key || resp.DifferentValues[0][1].Key != expectedResponse.DifferentValues[0][1].Key {
		t.Errorf("DifferentValues length mismatch: got %v want %v", len(resp.DifferentValues), len(expectedResponse.DifferentValues))
	}
	if len(resp.SameValues) != len(expectedResponse.SameValues) || resp.SameValues[0].Key != expectedResponse.SameValues[0].Key {
		t.Errorf("SameValues length mismatch: got %v want %v", len(resp.SameValues), len(expectedResponse.SameValues))
	}

}
