package apigee

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAppGroupAppKeyService_Create(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		appName      string
		input        *AppGroupAppKeyCreateRequest
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success with custom key",
			appGroupName: "app-group-1",
			appName:      "my-app",
			input: &AppGroupAppKeyCreateRequest{
				ConsumerKey:      "custom-key-123",
				ConsumerSecret:   "custom-secret-456",
				ExpiresInSeconds: 86400,
			},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupAppKey{
				ConsumerKey:    "custom-key-123",
				ConsumerSecret: "custom-secret-456",
				Status:         "approved",
			},
			wantErr: false,
		},
		{
			name:         "success without expiration",
			appGroupName: "app-group-1",
			appName:      "my-app",
			input: &AppGroupAppKeyCreateRequest{
				ConsumerKey:      "never-expires-key",
				ConsumerSecret:   "secret",
				ExpiresInSeconds: -1,
			},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupAppKey{
				ConsumerKey:    "never-expires-key",
				ConsumerSecret: "secret",
				Status:         "approved",
				ExpiresAt:      "-1",
			},
			wantErr: false,
		},
		{
			name:         "conflict - key already exists",
			appGroupName: "app-group-1",
			appName:      "my-app",
			input: &AppGroupAppKeyCreateRequest{
				ConsumerKey:    "existing-key",
				ConsumerSecret: "secret",
			},
			mockStatus:   http.StatusConflict,
			mockResponse: errorResponseBody(409, "Key already exists", "ALREADY_EXISTS"),
			wantErr:      true,
			errCheck:     IsConflict,
		},
		{
			name:         "app not found",
			appGroupName: "app-group-1",
			appName:      "nonexistent",
			input: &AppGroupAppKeyCreateRequest{
				ConsumerKey:    "key",
				ConsumerSecret: "secret",
			},
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "App not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "unauthorized",
			appGroupName: "app-group-1",
			appName:      "my-app",
			input: &AppGroupAppKeyCreateRequest{
				ConsumerKey:    "key",
				ConsumerSecret: "secret",
			},
			mockStatus:   http.StatusUnauthorized,
			mockResponse: errorResponseBody(401, "Invalid credentials", "UNAUTHENTICATED"),
			wantErr:      true,
			errCheck:     IsUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodPost,
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName + "/apps/" + tt.appName + "/keys",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroupAppKeys.Create(context.Background(), tt.appGroupName, tt.appName, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.ConsumerKey != tt.input.ConsumerKey {
				t.Errorf("result.ConsumerKey = %q, want %q", result.ConsumerKey, tt.input.ConsumerKey)
			}
		})
	}
}

func TestAppGroupAppKeyService_Get(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		appName      string
		consumerKey  string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "key-123",
			mockStatus:   http.StatusOK,
			mockResponse: &AppGroupAppKey{
				ConsumerKey:    "key-123",
				ConsumerSecret: "secret-456",
				Status:         "approved",
				APIProducts: []APIProductRef{
					{APIProduct: "product-1", Status: "approved"},
				},
			},
			wantErr: false,
		},
		{
			name:         "key not found",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Key not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "unauthorized",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "key-123",
			mockStatus:   http.StatusUnauthorized,
			mockResponse: errorResponseBody(401, "Invalid credentials", "UNAUTHENTICATED"),
			wantErr:      true,
			errCheck:     IsUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodGet,
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName + "/apps/" + tt.appName + "/keys/" + tt.consumerKey,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroupAppKeys.Get(context.Background(), tt.appGroupName, tt.appName, tt.consumerKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.ConsumerKey != tt.consumerKey {
				t.Errorf("result.ConsumerKey = %q, want %q", result.ConsumerKey, tt.consumerKey)
			}
		})
	}
}

func TestAppGroupAppKeyService_Update(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		appName      string
		consumerKey  string
		input        *AppGroupAppKey
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success - revoke key",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "key-123",
			input: &AppGroupAppKey{
				ConsumerKey: "key-123",
				Status:      "revoked",
			},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupAppKey{
				ConsumerKey: "key-123",
				Status:      "revoked",
			},
			wantErr: false,
		},
		{
			name:         "success - update API products",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "key-123",
			input: &AppGroupAppKey{
				ConsumerKey: "key-123",
				APIProducts: []APIProductRef{
					{APIProduct: "product-1", Status: "approved"},
					{APIProduct: "product-2", Status: "approved"},
				},
			},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupAppKey{
				ConsumerKey: "key-123",
				Status:      "approved",
				APIProducts: []APIProductRef{
					{APIProduct: "product-1", Status: "approved"},
					{APIProduct: "product-2", Status: "approved"},
				},
			},
			wantErr: false,
		},
		{
			name:         "key not found",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "nonexistent",
			input: &AppGroupAppKey{
				ConsumerKey: "nonexistent",
				Status:      "approved",
			},
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Key not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodPut,
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName + "/apps/" + tt.appName + "/keys/" + tt.consumerKey,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroupAppKeys.Update(context.Background(), tt.appGroupName, tt.appName, tt.consumerKey, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.Status != tt.input.Status && tt.input.Status != "" {
				t.Errorf("result.Status = %q, want %q", result.Status, tt.input.Status)
			}
		})
	}
}

func TestAppGroupAppKeyService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		appName      string
		consumerKey  string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "key-123",
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "key not found",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Key not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "forbidden",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "key-123",
			mockStatus:   http.StatusForbidden,
			mockResponse: errorResponseBody(403, "Access denied", "PERMISSION_DENIED"),
			wantErr:      true,
			errCheck:     IsForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodDelete,
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName + "/apps/" + tt.appName + "/keys/" + tt.consumerKey,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			err := client.AppGroupAppKeys.Delete(context.Background(), tt.appGroupName, tt.appName, tt.consumerKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
		})
	}
}

func TestAppGroupAppKeyService_UpdateAPIProductStatus(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		appName      string
		consumerKey  string
		apiProduct   string
		action       string
		mockStatus   int
		mockResponse interface{}
		wantPath     string
		wantQuery    string
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "approve product",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "key-123",
			apiProduct:   "product-1",
			action:       "approve",
			mockStatus:   http.StatusOK,
			wantPath:     "/organizations/test-org/appgroups/app-group-1/apps/my-app/keys/key-123/apiproducts/product-1",
			wantQuery:    "action=approve",
			wantErr:      false,
		},
		{
			name:         "revoke product",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "key-123",
			apiProduct:   "product-1",
			action:       "revoke",
			mockStatus:   http.StatusOK,
			wantPath:     "/organizations/test-org/appgroups/app-group-1/apps/my-app/keys/key-123/apiproducts/product-1",
			wantQuery:    "action=revoke",
			wantErr:      false,
		},
		{
			name:         "key not found",
			appGroupName: "app-group-1",
			appName:      "my-app",
			consumerKey:  "nonexistent",
			apiProduct:   "product-1",
			action:       "approve",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Key not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodPost,
				wantPath:   tt.wantPath,
				wantQuery:  tt.wantQuery,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			err := client.AppGroupAppKeys.UpdateAPIProductStatus(context.Background(), tt.appGroupName, tt.appName, tt.consumerKey, tt.apiProduct, tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAPIProductStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
		})
	}
}

func TestAppGroupAppKeyService_Generate(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		appName      string
		input        *AppGroupAppKeyGenerateRequest
		existingApp  *AppGroupApp
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			appName:      "my-app",
			input: &AppGroupAppKeyGenerateRequest{
				APIProducts:  []string{"product-1", "product-2"},
				KeyExpiresIn: 86400000,
			},
			existingApp: &AppGroupApp{
				Name:        "my-app",
				Status:      "approved",
				CallbackURL: "https://example.com/callback",
				Attributes: []Attribute{
					{Name: "env", Value: "prod"},
				},
				Scopes: []string{"read"},
			},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupApp{
				Name:        "my-app",
				Status:      "approved",
				CallbackURL: "https://example.com/callback",
				Credentials: []AppGroupAppKey{
					{
						ConsumerKey:    "generated-key",
						ConsumerSecret: "generated-secret",
						Status:         "approved",
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "app not found",
			appGroupName: "app-group-1",
			appName:      "nonexistent",
			input: &AppGroupAppKeyGenerateRequest{
				APIProducts: []string{"product-1"},
			},
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "App not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++

				// First call is GET to fetch existing app
				if callCount == 1 && r.Method == http.MethodGet {
					if tt.existingApp != nil {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(tt.existingApp)
					} else {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(tt.mockStatus)
						json.NewEncoder(w).Encode(tt.mockResponse)
					}
					return
				}

				// Second call is PUT to generate key
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				if tt.mockResponse != nil {
					json.NewEncoder(w).Encode(tt.mockResponse)
				}
			})

			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroupAppKeys.Generate(context.Background(), tt.appGroupName, tt.appName, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil {
				if len(result.Credentials) == 0 {
					t.Error("expected at least one credential in result")
				}
			}
		})
	}
}

func TestAppGroupAppKey_WithAttributes(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/appgroups/app-group-1/apps/my-app/keys/key-123",
		response: &AppGroupAppKey{
			ConsumerKey:    "key-123",
			ConsumerSecret: "secret-456",
			Status:         "approved",
			IssuedAt:       "1609459200000",
			ExpiresAt:      "1640995200000",
			Attributes: []Attribute{
				{Name: "env", Value: "production"},
				{Name: "team", Value: "api-team"},
			},
			APIProducts: []APIProductRef{
				{APIProduct: "product-1", Status: "approved"},
				{APIProduct: "product-2", Status: "pending"},
			},
		},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.AppGroupAppKeys.Get(context.Background(), "app-group-1", "my-app", "key-123")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if len(result.Attributes) != 2 {
		t.Errorf("Attributes count = %d, want 2", len(result.Attributes))
	}
	if len(result.APIProducts) != 2 {
		t.Errorf("APIProducts count = %d, want 2", len(result.APIProducts))
	}
	if result.IssuedAt != "1609459200000" {
		t.Errorf("IssuedAt = %q, want %q", result.IssuedAt, "1609459200000")
	}
	if result.ExpiresAt != "1640995200000" {
		t.Errorf("ExpiresAt = %q, want %q", result.ExpiresAt, "1640995200000")
	}
}

func TestAPIProductRef(t *testing.T) {
	ref := APIProductRef{
		APIProduct: "my-product",
		Status:     "approved",
	}

	if ref.APIProduct != "my-product" {
		t.Errorf("APIProduct = %q, want %q", ref.APIProduct, "my-product")
	}
	if ref.Status != "approved" {
		t.Errorf("Status = %q, want %q", ref.Status, "approved")
	}
}

func TestAppGroupAppKeyCreateRequest(t *testing.T) {
	req := &AppGroupAppKeyCreateRequest{
		ConsumerKey:      "custom-key",
		ConsumerSecret:   "custom-secret",
		ExpiresInSeconds: 3600,
	}

	if req.ConsumerKey != "custom-key" {
		t.Errorf("ConsumerKey = %q, want %q", req.ConsumerKey, "custom-key")
	}
	if req.ConsumerSecret != "custom-secret" {
		t.Errorf("ConsumerSecret = %q, want %q", req.ConsumerSecret, "custom-secret")
	}
	if req.ExpiresInSeconds != 3600 {
		t.Errorf("ExpiresInSeconds = %d, want %d", req.ExpiresInSeconds, 3600)
	}
}

func TestAppGroupAppKeyGenerateRequest(t *testing.T) {
	req := &AppGroupAppKeyGenerateRequest{
		APIProducts:  []string{"product-1", "product-2"},
		KeyExpiresIn: 86400000,
	}

	if len(req.APIProducts) != 2 {
		t.Errorf("APIProducts count = %d, want 2", len(req.APIProducts))
	}
	if req.KeyExpiresIn != 86400000 {
		t.Errorf("KeyExpiresIn = %d, want %d", req.KeyExpiresIn, 86400000)
	}
}
