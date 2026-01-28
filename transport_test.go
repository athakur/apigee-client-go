package apigee

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_buildPath(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		segments []string
		want     string
	}{
		{
			name:     "simple path",
			baseURL:  "https://api.example.com/v1",
			segments: []string{"organizations", "my-org"},
			want:     "https://api.example.com/v1/organizations/my-org",
		},
		{
			name:     "path with special characters",
			baseURL:  "https://api.example.com/v1",
			segments: []string{"organizations", "my org", "apps"},
			want:     "https://api.example.com/v1/organizations/my%20org/apps",
		},
		{
			name:     "path with slashes in segment",
			baseURL:  "https://api.example.com/v1",
			segments: []string{"organizations", "org/name"},
			want:     "https://api.example.com/v1/organizations/org%2Fname",
		},
		{
			name:     "empty segments",
			baseURL:  "https://api.example.com/v1",
			segments: []string{},
			want:     "https://api.example.com/v1/",
		},
		{
			name:     "single segment",
			baseURL:  "https://api.example.com/v1",
			segments: []string{"organizations"},
			want:     "https://api.example.com/v1/organizations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{BaseURL: tt.baseURL}
			got := c.buildPath(tt.segments...)
			if got != tt.want {
				t.Errorf("buildPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAddQueryParams(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		opts    *ListOptions
		want    string
		wantErr bool
	}{
		{
			name:   "nil options",
			urlStr: "https://api.example.com/v1/resources",
			opts:   nil,
			want:   "https://api.example.com/v1/resources",
		},
		{
			name:   "empty options",
			urlStr: "https://api.example.com/v1/resources",
			opts:   &ListOptions{},
			want:   "https://api.example.com/v1/resources",
		},
		{
			name:   "page size only",
			urlStr: "https://api.example.com/v1/resources",
			opts:   &ListOptions{PageSize: 10},
			want:   "https://api.example.com/v1/resources?pageSize=10",
		},
		{
			name:   "page token only",
			urlStr: "https://api.example.com/v1/resources",
			opts:   &ListOptions{PageToken: "abc123"},
			want:   "https://api.example.com/v1/resources?pageToken=abc123",
		},
		{
			name:   "filter only",
			urlStr: "https://api.example.com/v1/resources",
			opts:   &ListOptions{Filter: "status=active"},
			want:   "https://api.example.com/v1/resources?filter=status%3Dactive",
		},
		{
			name:   "all options",
			urlStr: "https://api.example.com/v1/resources",
			opts:   &ListOptions{PageSize: 25, PageToken: "token123", Filter: "name=test"},
			// Note: URL encoding order may vary, so we check for presence
			want: "https://api.example.com/v1/resources?filter=name%3Dtest&pageSize=25&pageToken=token123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addQueryParams(tt.urlStr, tt.opts)
			if got != tt.want {
				t.Errorf("addQueryParams() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClient_newRequest(t *testing.T) {
	c := &Client{BaseURL: "https://api.example.com/v1"}
	ctx := context.Background()

	tests := []struct {
		name            string
		method          string
		url             string
		body            interface{}
		wantContentType string
		wantAccept      string
		wantErr         bool
	}{
		{
			name:       "GET request without body",
			method:     http.MethodGet,
			url:        "https://api.example.com/v1/resources",
			body:       nil,
			wantAccept: "application/json",
		},
		{
			name:            "POST request with body",
			method:          http.MethodPost,
			url:             "https://api.example.com/v1/resources",
			body:            map[string]string{"name": "test"},
			wantContentType: "application/json",
			wantAccept:      "application/json",
		},
		{
			name:            "PUT request with body",
			method:          http.MethodPut,
			url:             "https://api.example.com/v1/resources/1",
			body:            map[string]string{"name": "updated"},
			wantContentType: "application/json",
			wantAccept:      "application/json",
		},
		{
			name:       "DELETE request without body",
			method:     http.MethodDelete,
			url:        "https://api.example.com/v1/resources/1",
			body:       nil,
			wantAccept: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := c.newRequest(ctx, tt.method, tt.url, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("newRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if req.Method != tt.method {
				t.Errorf("method = %q, want %q", req.Method, tt.method)
			}
			if req.URL.String() != tt.url {
				t.Errorf("url = %q, want %q", req.URL.String(), tt.url)
			}
			if got := req.Header.Get("Accept"); got != tt.wantAccept {
				t.Errorf("Accept header = %q, want %q", got, tt.wantAccept)
			}
			if tt.body != nil {
				if got := req.Header.Get("Content-Type"); got != tt.wantContentType {
					t.Errorf("Content-Type header = %q, want %q", got, tt.wantContentType)
				}
			}
		})
	}
}

func TestClient_newRequest_marshalError(t *testing.T) {
	c := &Client{BaseURL: "https://api.example.com/v1"}
	ctx := context.Background()

	// Channel cannot be marshaled to JSON
	body := make(chan int)
	_, err := c.newRequest(ctx, http.MethodPost, "https://api.example.com/v1/resources", body)
	if err == nil {
		t.Error("expected error for unmarshalable body, got nil")
	}
	if !strings.Contains(err.Error(), "marshal") {
		t.Errorf("error should mention marshal, got: %v", err)
	}
}

func TestClient_doRequest_authHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Authorization header = %q, want %q", auth, "Bearer test-token")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	client, server := setupTestClient(t, handler)
	req, _ := client.newRequest(context.Background(), http.MethodGet, server.URL+"/test", nil)
	_, err := client.doRequest(req)
	if err != nil {
		t.Errorf("doRequest() error = %v", err)
	}
}

func TestClient_doRequest_tokenError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(context.Background(), "test-org",
		WithBaseURL(server.URL),
		WithTokenSource(&mockTokenSource{err: errors.New("token error")}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	req, _ := client.newRequest(context.Background(), http.MethodGet, server.URL+"/test", nil)
	_, err = client.doRequest(req)
	if err == nil {
		t.Error("expected error for token failure, got nil")
	}
	if !strings.Contains(err.Error(), "token") {
		t.Errorf("error should mention token, got: %v", err)
	}
}

func TestClient_doRequest_httpError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   interface{}
		errCheck   func(error) bool
	}{
		{
			name:       "401 unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   errorResponseBody(401, "Invalid credentials", "UNAUTHENTICATED"),
			errCheck:   IsUnauthorized,
		},
		{
			name:       "403 forbidden",
			statusCode: http.StatusForbidden,
			response:   errorResponseBody(403, "Access denied", "PERMISSION_DENIED"),
			errCheck:   IsForbidden,
		},
		{
			name:       "404 not found",
			statusCode: http.StatusNotFound,
			response:   errorResponseBody(404, "Resource not found", "NOT_FOUND"),
			errCheck:   IsNotFound,
		},
		{
			name:       "409 conflict",
			statusCode: http.StatusConflict,
			response:   errorResponseBody(409, "Resource already exists", "ALREADY_EXISTS"),
			errCheck:   IsConflict,
		},
		{
			name:       "500 internal server error",
			statusCode: http.StatusInternalServerError,
			response:   errorResponseBody(500, "Internal error", "INTERNAL"),
			errCheck:   func(err error) bool { return err != nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := setupTestClient(t, jsonHandler(t, tt.statusCode, tt.response))
			req, _ := client.newRequest(context.Background(), http.MethodGet, client.BaseURL+"/test", nil)
			_, err := client.doRequest(req)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
		})
	}
}

func TestClient_do(t *testing.T) {
	type testResponse struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name       string
		method     string
		statusCode int
		response   interface{}
		wantResult *testResponse
		wantErr    bool
	}{
		{
			name:       "successful GET",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			response:   &testResponse{Name: "test", Value: 42},
			wantResult: &testResponse{Name: "test", Value: 42},
		},
		{
			name:       "successful POST",
			method:     http.MethodPost,
			statusCode: http.StatusCreated,
			response:   &testResponse{Name: "created", Value: 1},
			wantResult: &testResponse{Name: "created", Value: 1},
		},
		{
			name:       "error response",
			method:     http.MethodGet,
			statusCode: http.StatusNotFound,
			response:   errorResponseBody(404, "Not found", "NOT_FOUND"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := setupTestClient(t, jsonHandler(t, tt.statusCode, tt.response))

			var result testResponse
			err := client.do(context.Background(), tt.method, client.BaseURL+"/test", nil, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantResult != nil {
				if result.Name != tt.wantResult.Name || result.Value != tt.wantResult.Value {
					t.Errorf("do() result = %+v, want %+v", result, tt.wantResult)
				}
			}
		})
	}
}

func TestClient_do_emptyResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	client, server := setupTestClient(t, handler)
	err := client.do(context.Background(), http.MethodDelete, server.URL+"/test", nil, nil)
	if err != nil {
		t.Errorf("do() error = %v, want nil", err)
	}
}

func TestClient_do_requestBody(t *testing.T) {
	type requestBody struct {
		Name string `json:"name"`
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}

		var got requestBody
		if err := json.Unmarshal(body, &got); err != nil {
			t.Errorf("failed to unmarshal request body: %v", err)
		}

		if got.Name != "test-name" {
			t.Errorf("request body name = %q, want %q", got.Name, "test-name")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	client, server := setupTestClient(t, handler)
	reqBody := &requestBody{Name: "test-name"}
	err := client.do(context.Background(), http.MethodPost, server.URL+"/test", reqBody, nil)
	if err != nil {
		t.Errorf("do() error = %v", err)
	}
}

func TestClient_do_unmarshalError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	})

	client, server := setupTestClient(t, handler)

	type response struct {
		Name string `json:"name"`
	}
	var result response
	err := client.do(context.Background(), http.MethodGet, server.URL+"/test", nil, &result)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "unmarshal") {
		t.Errorf("error should mention unmarshal, got: %v", err)
	}
}
