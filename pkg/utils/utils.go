package utils

import (
	"encoding/json"
	"net/http"
)

type SuccessMessage[T any] struct {
	Status  int  `json:"status"`
	Data    T    `json:"data,omitempty"`
	Success bool `json:"success"`
}

func RespondMessage(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(newSuccessMessage(status, data))
}

func newSuccessMessage[T any](status int, data T) *SuccessMessage[T] {
	success := true
	if status >= http.StatusBadRequest {
		success = false
	}

	return &SuccessMessage[T]{
		Status:  status,
		Data:    data,
		Success: success,
	}
}
