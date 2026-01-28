package apigee

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "with message",
			err: &Error{
				StatusCode: 404,
				Code:       5,
				Message:    "Resource not found",
			},
			want: "apigee: Resource not found (status 404, code 5)",
		},
		{
			name: "without message",
			err: &Error{
				StatusCode: 500,
			},
			want: "apigee: HTTP 500",
		},
		{
			name: "with all fields",
			err: &Error{
				StatusCode: 409,
				Code:       6,
				Message:    "Already exists",
				Status:     "ALREADY_EXISTS",
			},
			want: "apigee: Already exists (status 409, code 6)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		body           string
		wantStatusCode int
		wantCode       int
		wantMessage    string
		wantStatus     string
	}{
		{
			name:       "valid error response",
			statusCode: 404,
			body: `{
				"error": {
					"code": 404,
					"message": "Resource not found",
					"status": "NOT_FOUND"
				}
			}`,
			wantStatusCode: 404,
			wantCode:       404,
			wantMessage:    "Resource not found",
			wantStatus:     "NOT_FOUND",
		},
		{
			name:           "empty body",
			statusCode:     500,
			body:           "",
			wantStatusCode: 500,
		},
		{
			name:           "invalid JSON",
			statusCode:     502,
			body:           "Bad Gateway",
			wantStatusCode: 502,
		},
		{
			name:           "empty error object",
			statusCode:     503,
			body:           `{"error": {}}`,
			wantStatusCode: 503,
		},
		{
			name:           "no error field",
			statusCode:     400,
			body:           `{"message": "Bad request"}`,
			wantStatusCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			err := parseError(resp, []byte(tt.body))

			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("parseError() returned non-Error type: %T", err)
			}

			if apiErr.StatusCode != tt.wantStatusCode {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.wantStatusCode)
			}
			if apiErr.Code != tt.wantCode {
				t.Errorf("Code = %d, want %d", apiErr.Code, tt.wantCode)
			}
			if apiErr.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", apiErr.Message, tt.wantMessage)
			}
			if apiErr.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", apiErr.Status, tt.wantStatus)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "404 error",
			err:  &Error{StatusCode: http.StatusNotFound},
			want: true,
		},
		{
			name: "wrapped 404 error",
			err:  fmt.Errorf("wrapped: %w", &Error{StatusCode: http.StatusNotFound}),
			want: true,
		},
		{
			name: "400 error",
			err:  &Error{StatusCode: http.StatusBadRequest},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsConflict(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "409 error",
			err:  &Error{StatusCode: http.StatusConflict},
			want: true,
		},
		{
			name: "wrapped 409 error",
			err:  fmt.Errorf("wrapped: %w", &Error{StatusCode: http.StatusConflict}),
			want: true,
		},
		{
			name: "400 error",
			err:  &Error{StatusCode: http.StatusBadRequest},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsConflict(tt.err); got != tt.want {
				t.Errorf("IsConflict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsForbidden(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "403 error",
			err:  &Error{StatusCode: http.StatusForbidden},
			want: true,
		},
		{
			name: "wrapped 403 error",
			err:  fmt.Errorf("wrapped: %w", &Error{StatusCode: http.StatusForbidden}),
			want: true,
		},
		{
			name: "401 error",
			err:  &Error{StatusCode: http.StatusUnauthorized},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsForbidden(tt.err); got != tt.want {
				t.Errorf("IsForbidden() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "401 error",
			err:  &Error{StatusCode: http.StatusUnauthorized},
			want: true,
		},
		{
			name: "wrapped 401 error",
			err:  fmt.Errorf("wrapped: %w", &Error{StatusCode: http.StatusUnauthorized}),
			want: true,
		},
		{
			name: "403 error",
			err:  &Error{StatusCode: http.StatusForbidden},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnauthorized(tt.err); got != tt.want {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorAsInterface(t *testing.T) {
	// Verify that *Error implements the error interface
	var _ error = (*Error)(nil)

	// Verify errors.As works with wrapped errors
	originalErr := &Error{StatusCode: 404, Message: "not found"}
	wrappedErr := fmt.Errorf("operation failed: %w", originalErr)

	var apiErr *Error
	if !errors.As(wrappedErr, &apiErr) {
		t.Error("errors.As() should find *Error in wrapped error")
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
}
