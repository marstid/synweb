package errors

import (
	"errors"
	"fmt"
)

var (
	ErrMissingQuery    = New("MISSING_QUERY", "Query parameter is required")
	ErrInvalidQuery    = New("INVALID_QUERY", "Query must be a non-empty string")
	ErrAPIKeyMissing   = New("API_KEY_MISSING", "SYNTHETIC_API_KEY environment variable is required")
	ErrAPIKeyInvalid   = New("API_KEY_INVALID", "Invalid API key format")
	ErrCircuitOpen     = New("CIRCUIT_OPEN", "Circuit breaker is open, request rejected")
	ErrTooManyRetries  = New("TOO_MANY_RETRIES", "Maximum retry attempts exceeded")
	ErrRequestTimeout  = New("REQUEST_TIMEOUT", "Request timed out")
	ErrNetworkError    = New("NETWORK_ERROR", "Network connection failed")
	ErrServerError     = New("SERVER_ERROR", "Server returned an error")
	ErrInvalidResponse = New("INVALID_RESPONSE", "Invalid response format from server")
	ErrRateLimited     = New("RATE_LIMITED", "Rate limit exceeded")
)

type Error struct {
	Code    string
	Message string
	Err     error
}

func New(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) With(err error) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %d %s - %s", e.StatusCode, e.Status, e.Body)
}

func NewAPIError(statusCode int, status, body string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Status:     status,
		Body:       body,
	}
}
