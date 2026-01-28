package apigee

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Error represents an error response from the Apigee API.
type Error struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int `json:"-"`
	// Code is the error code returned by the API.
	Code int `json:"code,omitempty"`
	// Message is the error message returned by the API.
	Message string `json:"message,omitempty"`
	// Status is the status string returned by the API.
	Status string `json:"status,omitempty"`
}

// errorResponse is the wrapper for API error responses.
type errorResponse struct {
	Error *Error `json:"error,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("apigee: %s (status %d, code %d)", e.Message, e.StatusCode, e.Code)
	}
	return fmt.Sprintf("apigee: HTTP %d", e.StatusCode)
}

// parseError parses an error response from the API.
func parseError(resp *http.Response, body []byte) error {
	apiErr := &Error{
		StatusCode: resp.StatusCode,
	}

	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != nil {
		apiErr.Code = errResp.Error.Code
		apiErr.Message = errResp.Error.Message
		apiErr.Status = errResp.Error.Status
	}

	return apiErr
}

// IsNotFound reports whether err is a 404 Not Found error.
func IsNotFound(err error) bool {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsConflict reports whether err is a 409 Conflict error.
func IsConflict(err error) bool {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusConflict
	}
	return false
}

// IsForbidden reports whether err is a 403 Forbidden error.
func IsForbidden(err error) bool {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// IsUnauthorized reports whether err is a 401 Unauthorized error.
func IsUnauthorized(err error) bool {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}
