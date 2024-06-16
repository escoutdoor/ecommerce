package respond

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
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

func ValidationError(errors validator.ValidationErrors) error {
	var msgs []string

	for _, err := range errors {
		switch err.ActualTag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "min":
			if err.Kind() == reflect.String {
				msgs = append(msgs, fmt.Sprintf("field %s should be at least %v characters long", err.Field(), err.Param()))
				break
			}

			msgs = append(msgs, fmt.Sprintf("field %s should contain at least %v items", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("field %s should not exceed %s symbols long", err.Field(), err.Param()))
		case "containsany":
			msgs = append(msgs, fmt.Sprintf("field %s should contain at least one special character (%s)", err.Field(), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return fmt.Errorf(strings.Join(msgs, ", "))
}
