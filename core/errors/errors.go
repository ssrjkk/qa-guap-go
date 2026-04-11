package errors

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Message    string
	URL        string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error [%d]: %s (url: %s)", e.StatusCode, e.Message, e.URL)
}

func NewAPIError(statusCode int, message, url string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		URL:        url,
	}
}

func HandleHTTPError(resp *http.Response, url string) error {
	if resp.StatusCode >= 400 {
		return NewAPIError(resp.StatusCode, fmt.Sprintf("HTTP %d", resp.StatusCode), url)
	}
	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s", e.Field, e.Message)
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("Retryable error: %v", e.Err)
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

func NewRetryableError(err error) *RetryableError {
	return &RetryableError{Err: err}
}

func IsRetryable(err error) bool {
	_, ok := err.(*RetryableError)
	return ok
}
