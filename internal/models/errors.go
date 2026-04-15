// Package models contains shared data structures and error types for DDGS.
package models

import "fmt"

// DDGSError represents a base error type for DDGS operations.
type DDGSError struct {
	Message string
}

func (e *DDGSError) Error() string {
	return e.Message
}

// NewDDGSError creates a new DDGSError.
func NewDDGSError(msg string) *DDGSError {
	return &DDGSError{Message: msg}
}

// RateLimitError represents a rate limit exceeded error.
type RateLimitError struct {
	DDGSError
}

// NewRateLimitError creates a new RateLimitError.
func NewRateLimitError(msg string) *RateLimitError {
	return &RateLimitError{DDGSError: DDGSError{Message: msg}}
}

// TimeoutError represents a timeout error.
type TimeoutError struct {
	DDGSError
	Err error
}

func (e *TimeoutError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("timeout: %v", e.Err)
	}
	return e.Message
}

// NewTimeoutError creates a new TimeoutError.
func NewTimeoutError(err error) *TimeoutError {
	return &TimeoutError{
		DDGSError: DDGSError{Message: "request timed out"},
		Err:       err,
	}
}

// IsTimeoutError checks if an error is a TimeoutError.
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*TimeoutError)
	return ok
}
