package response

import (
	"encoding/json"
	"net/http"
)

// Success define the structure of success for http response
type Success struct {
	statusCode int
	result     interface{}
}

// NewSuccess create new Success
func NewSuccess(status int, result interface{}) Success {
	return Success{
		statusCode: status,
		result:     result,
	}
}

// Send return a response with JSON format
func (r Success) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.statusCode)
	return json.NewEncoder(w).Encode(r.result)
}
