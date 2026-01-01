package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func ParseJSON(r *http.Request, payload any) error {

	if r.Body == nil {
		log.Fatal("Need Body in the http response")
	}

	return json.NewDecoder(r.Body).Decode(payload)

}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	WriteJSON(w, statusCode, map[string]string{"error": err.Error()})
}
