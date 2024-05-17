package respond

import (
	"encoding/json"
	"net/http"
)

type ApiError struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ApiError{Error: err.Error()}
	json.NewEncoder(w).Encode(resp)
}
