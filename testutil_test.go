package apigee

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

// mockTokenSource is a mock OAuth2 token source for testing.
type mockTokenSource struct {
	token *oauth2.Token
	err   error
}

func (m *mockTokenSource) Token() (*oauth2.Token, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.token != nil {
		return m.token, nil
	}
	return &oauth2.Token{AccessToken: "test-token"}, nil
}

// setupTestClient creates a client pointing to a test server.
func setupTestClient(t *testing.T, handler http.Handler) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := NewClient(context.Background(), "test-org",
		WithBaseURL(server.URL),
		WithTokenSource(&mockTokenSource{}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client, server
}

// jsonHandler creates a handler that returns JSON with the specified status code.
func jsonHandler(t *testing.T, statusCode int, body interface{}) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != nil {
			if err := json.NewEncoder(w).Encode(body); err != nil {
				t.Errorf("failed to encode response body: %v", err)
			}
		}
	}
}

// errorResponse creates a standard Apigee error response.
func errorResponseBody(code int, message, status string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"status":  status,
		},
	}
}

// requestValidator is a handler that validates request properties.
type requestValidator struct {
	t              *testing.T
	wantMethod     string
	wantPath       string
	wantQuery      string
	wantAuthHeader string
	wantBody       interface{}
	response       interface{}
	statusCode     int
}

func (rv *requestValidator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rv.wantMethod != "" && r.Method != rv.wantMethod {
		rv.t.Errorf("method = %q, want %q", r.Method, rv.wantMethod)
	}
	if rv.wantPath != "" && r.URL.Path != rv.wantPath {
		rv.t.Errorf("path = %q, want %q", r.URL.Path, rv.wantPath)
	}
	if rv.wantQuery != "" && r.URL.RawQuery != rv.wantQuery {
		rv.t.Errorf("query = %q, want %q", r.URL.RawQuery, rv.wantQuery)
	}
	if rv.wantAuthHeader != "" {
		gotAuth := r.Header.Get("Authorization")
		if gotAuth != rv.wantAuthHeader {
			rv.t.Errorf("Authorization header = %q, want %q", gotAuth, rv.wantAuthHeader)
		}
	}

	statusCode := rv.statusCode
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if rv.response != nil {
		if err := json.NewEncoder(w).Encode(rv.response); err != nil {
			rv.t.Errorf("failed to encode response: %v", err)
		}
	}
}

// muxHandler creates a multiplexer handler for multiple endpoints.
func muxHandler(handlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + " " + r.URL.Path
		if handler, ok := handlers[key]; ok {
			handler(w, r)
			return
		}
		http.NotFound(w, r)
	}
}
