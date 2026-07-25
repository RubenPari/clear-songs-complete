// Package shared defines domain-level interfaces, error sentinels, and cache key
// constants used across the application. It establishes the contracts that
// infrastructure implementations must satisfy.
package shared

import "errors"

var (
	// ErrNotFound indicates that a requested resource was not found.
	ErrNotFound = errors.New("resource not found")
	
	// ErrUnauthorized indicates that the user is not authenticated or lacks permissions.
	ErrUnauthorized = errors.New("unauthorized access")
	
	// ErrValidation indicates that the provided input is invalid.
	ErrValidation = errors.New("validation failed")
	
	// ErrInternal indicates an unexpected internal server error.
	ErrInternal = errors.New("internal server error")
	
	// ErrExternalAPI indicates a failure when communicating with an external service (e.g., Spotify).
	ErrExternalAPI = errors.New("external API error")
)
