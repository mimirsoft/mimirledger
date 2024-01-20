package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// get an object, marshal and return it as JSON
func RespondOK(w http.ResponseWriter, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("failed json encode response payload: %w", err)
	}
	return nil
}

// respond with no content and 204 header
func RespondNoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// get an object, marshal and return it as JSON
func RespondRedirectStatusFound(w http.ResponseWriter, location string) error {
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusFound)
	return nil
}
