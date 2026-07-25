package dto

import (
	"time"
)

// APIResponse represents a standardized API response using Go Generics
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Error represents an error in the API response
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta contains metadata about the response
type Meta struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

// NewSuccess creates a successful response with typed data
func NewSuccess[T any](data T) APIResponse[T] {
	return APIResponse[T]{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// NewError creates a generic error response with the given machine-readable code
// and human-readable message.
func NewError(code, message string) APIResponse[any] {
	return APIResponse[any]{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
		},
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// ValidationErr creates a 400-style error response with the VALIDATION_ERROR code.
func ValidationErr(message string) APIResponse[any] {
	return NewError("VALIDATION_ERROR", message)
}

// InternalErr creates a 500-style error response with the INTERNAL_ERROR code.
func InternalErr(message string) APIResponse[any] {
	return NewError("INTERNAL_ERROR", message)
}

// NotFoundErr creates a 404-style error response with the NOT_FOUND code and a
// message indicating the named resource was not found.
func NotFoundErr(resource string) APIResponse[any] {
	return NewError("NOT_FOUND", resource+" not found")
}

// UnauthorizedErr creates a 401-style error response with the UNAUTHORIZED code
// and a fixed message indicating authentication is required.
func UnauthorizedErr() APIResponse[any] {
	return NewError("UNAUTHORIZED", "Authentication required")
}
