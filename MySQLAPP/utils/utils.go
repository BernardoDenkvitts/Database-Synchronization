package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, response any, header *map[string]string) error {
	w.Header().Add("Content-Type", "application-json")
	if header != nil && len(*header) > 0 {
		for key, value := range *header {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}
