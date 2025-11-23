package app_err

import (
	"time"

	"github.com/google/uuid"
)

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	RequestID uuid.UUID     `json:"request_id"`        // Уникальный ID запроса
	Timestamp time.Time     `json:"timestamp"`         // Время ошибки
	Error     string        `json:"error"`             // Общее описание
	Details   []ErrorDetail `json:"details,omitempty"` // Детальные ошибки
	Path      string        `json:"path"`              // URL путь
}

func WriteErrorWithDetails(status int, message, path string, details []ErrorDetail) *ErrorResponse {
	return &ErrorResponse{
		RequestID: uuid.New(),
		Timestamp: time.Now().UTC(),
		Error:     message,
		Details:   details,
		Path:      path,
	}
}
func WriteError(message, path string) *ErrorResponse {
	return &ErrorResponse{
		RequestID: uuid.New(),
		Timestamp: time.Now().UTC(),
		Error:     message,
		Path:      path,
	}
}
