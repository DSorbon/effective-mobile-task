package response

import (
	"encoding/json"
	"net/http"

	"github.com/DSorbon/effective-mobile-task/internal/models"
)

type ResponseMessage struct {
	Message string `json:"message,omitempty"`
}

type ResponseValidationErrors struct {
	Message        string          `json:"message"`
	ValidateErrors json.RawMessage `json:"validate_errors" swaggertype:"object,string"`
}

type ResponseData struct {
	Data any `json:"data"`
}

type ResponsePagintion struct {
	Data any               `json:"data"`
	Page models.Pagination `json:"page"`
}

func WithBody(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

func WithoutBody(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

// ValidateErrors sends validation errors as a JSON response
func ValidateErrors(w http.ResponseWriter, statusCode int, errors []byte) {
	response := ResponseValidationErrors{
		Message:        "validation failed",
		ValidateErrors: json.RawMessage(errors),
	}
	WithBody(w, statusCode, response)
}
