package edlib

import (
	"errors"
	"fmt"
)

// Sentinel errors for common failure cases
var (
	// ErrInvalidAlgorithm is returned when an unknown or unsupported algorithm is specified
	ErrInvalidAlgorithm = errors.New("invalid or unsupported algorithm")

	// ErrInvalidConfiguration is returned when the configuration is invalid
	ErrInvalidConfiguration = errors.New("invalid configuration")

	// ErrEmptyStrings is returned when one or both input strings are empty where it's not allowed
	ErrEmptyStrings = errors.New("input strings cannot be empty")

	// ErrUnequalLength is returned when strings must be of equal length but are not
	ErrUnequalLength = errors.New("strings must be of equal length")

	// ErrInvalidThreshold is returned when threshold is out of valid range [0, 1]
	ErrInvalidThreshold = errors.New("threshold must be between 0 and 1")

	// ErrInvalidWeights is returned when cost weights are invalid (e.g., negative values)
	ErrInvalidWeights = errors.New("invalid cost weights")

	// ErrContextCanceled is returned when the context is canceled during computation
	ErrContextCanceled = errors.New("operation canceled")

	// ErrContextTimeout is returned when the context times out during computation
	ErrContextTimeout = errors.New("operation timed out")

	// ErrNoResults is returned when no results match the search criteria
	ErrNoResults = errors.New("no results found")

	// ErrInvalidSearchParams is returned when search parameters are invalid
	ErrInvalidSearchParams = errors.New("invalid search parameters")
)

// ConfigurationError wraps configuration-related errors with additional context
type ConfigurationError struct {
	Field   string
	Value   interface{}
	Reason  string
	wrapped error
}

func (e *ConfigurationError) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("configuration error [%s=%v]: %s: %v", e.Field, e.Value, e.Reason, e.wrapped)
	}
	return fmt.Sprintf("configuration error [%s=%v]: %s", e.Field, e.Value, e.Reason)
}

func (e *ConfigurationError) Unwrap() error {
	return e.wrapped
}

// NewConfigurationError creates a new configuration error
func NewConfigurationError(field string, value interface{}, reason string) error {
	return &ConfigurationError{
		Field:   field,
		Value:   value,
		Reason:  reason,
		wrapped: ErrInvalidConfiguration,
	}
}

// AlgorithmError wraps algorithm-specific errors with additional context
type AlgorithmError struct {
	Algorithm string
	Operation string
	Reason    string
	wrapped   error
}

func (e *AlgorithmError) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("algorithm error [%s/%s]: %s: %v", e.Algorithm, e.Operation, e.Reason, e.wrapped)
	}
	return fmt.Sprintf("algorithm error [%s/%s]: %s", e.Algorithm, e.Operation, e.Reason)
}

func (e *AlgorithmError) Unwrap() error {
	return e.wrapped
}

// NewAlgorithmError creates a new algorithm error
func NewAlgorithmError(algorithm, operation, reason string, wrapped error) error {
	return &AlgorithmError{
		Algorithm: algorithm,
		Operation: operation,
		Reason:    reason,
		wrapped:   wrapped,
	}
}
