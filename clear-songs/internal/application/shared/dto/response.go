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

// NewError creates an error response
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

// ValidationErr creates a validation error response
func ValidationErr(message string) APIResponse[any] {
	return NewError("VALIDATION_ERROR", message)
}

// InternalErr creates an internal server error response
func InternalErr(message string) APIResponse[any] {
	return NewError("INTERNAL_ERROR", message)
}

// NotFoundErr creates a not found error response
func NotFoundErr(resource string) APIResponse[any] {
	return NewError("NOT_FOUND", resource+" not found")
}

// UnauthorizedErr creates an unauthorized error response
func UnauthorizedErr() APIResponse[any] {
	return NewError("UNAUTHORIZED", "Authentication required")
}
