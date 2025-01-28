package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// Apply MaxBytesReader to ensure size limit is enforced
	maxBytes := 1_048_578 // 1MB limit
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Create a JSON decoder and apply settings
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow fields not defined in the struct

	// Attempt to decode the JSON body
	if err := decoder.Decode(data); err != nil {
		// Handle decoding errors
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding JSON: %s", err.Error()))
		return err
	}

	return nil
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}
