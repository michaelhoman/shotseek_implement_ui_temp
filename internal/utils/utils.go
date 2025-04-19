package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// Apply MaxBytesReader to ensure size limit is enforced
	maxBytes := 1_048_578 // 1MB limit
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Create a JSON decoder and apply settings
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow fields not defined in the struct

	// Attempt to decode the JSON body
	if err := decoder.Decode(data); err != nil {
		// Handle decoding errors
		WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding JSON: %s", err.Error()))
		return err
	}

	return nil
}

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return WriteJSON(w, status, &envelope{Error: message})
}

func JsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return WriteJSON(w, status, &envelope{Data: data})
}

// func (app *application) writeStringJSON(w http.ResponseWriter, status int, message string) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)

// 	// Send the message as raw text inside a JSON structure
// 	jsonResponse := fmt.Sprintf("{\"message\":\"%s\"}", message)

// 	_, err := w.Write([]byte(jsonResponse))
// 	return err
// }

func WriteMessagePlain(w http.ResponseWriter, status int, message string) error {
	w.Header().Set("Content-Type", "text/plain") // Change to text/html if needed
	w.WriteHeader(status)

	_, err := w.Write([]byte(message))
	return err
}

func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}
