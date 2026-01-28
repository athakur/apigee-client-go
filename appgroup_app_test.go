package apigee

import (
	"context"
	"net/http"
	"testing"
)

func TestAppGroupAppService_Create(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		input        *AppGroupApp
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			input: &AppGroupApp{
				Name:        "my-app",
				Status:      "approved",
				CallbackURL: "https://example.com/callback",
				APIProducts: []string{"product-1", "product-2"},
			},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupApp{
				Name:        "my-app",
				AppID:       "generated-app-id",
				AppGroup:    "app-group-1",
				Status:      "approved",
				CallbackURL: "https://example.com/callback",
				APIProducts: []string{"product-1", "product-2"},
			},
			wantErr: false,
		},
		{
			name:         "conflict",
			appGroupName: "app-group-1",
			input:        &AppGroupApp{Name: "existing-app"},
			mockStatus:   http.StatusConflict,
			mockResponse: errorResponseBody(409, "App already exists", "ALREADY_EXISTS"),
			wantErr:      true,
			errCheck:     IsConflict,
		},
		{
			name:         "app group not found",
			appGroupName: "nonexistent",
			input:        &AppGroupApp{Name: "my-app"},
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "App group not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "unauthorized",
			appGroupName: "app-group-1",
			input:        &AppGroupApp{Name: "my-app"},
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
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName + "/apps",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroupApps.Create(context.Background(), tt.appGroupName, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.Name != tt.input.Name {
				t.Errorf("result.Name = %q, want %q", result.Name, tt.input.Name)
			}
		})
	}
}

func TestAppGroupAppService_Get(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		appName      string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			appName:      "my-app",
			mockStatus:   http.StatusOK,
			mockResponse: &AppGroupApp{
				Name:     "my-app",
				AppID:    "app-id-123",
				AppGroup: "app-group-1",
				Status:   "approved",
				Credentials: []AppGroupAppKey{
					{ConsumerKey: "key-1", Status: "approved"},
				},
			},
			wantErr: false,
		},
		{
			name:         "app not found",
			appGroupName: "app-group-1",
			appName:      "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "App not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "unauthorized",
			appGroupName: "app-group-1",
			appName:      "my-app",
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
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName + "/apps/" + tt.appName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroupApps.Get(context.Background(), tt.appGroupName, tt.appName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.Name != tt.appName {
				t.Errorf("result.Name = %q, want %q", result.Name, tt.appName)
			}
		})
	}
}

func TestAppGroupAppService_Update(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		appName      string
		input        *AppGroupApp
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			appName:      "my-app",
			input: &AppGroupApp{
				Name:        "my-app",
				Status:      "revoked",
				CallbackURL: "https://updated.example.com/callback",
			},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupApp{
				Name:        "my-app",
				Status:      "revoked",
				CallbackURL: "https://updated.example.com/callback",
			},
			wantErr: false,
		},
		{
			name:         "app not found",
			appGroupName: "app-group-1",
			appName:      "nonexistent",
			input: &AppGroupApp{
				Name:   "nonexistent",
				Status: "approved",
			},
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "App not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodPut,
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName + "/apps/" + tt.appName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroupApps.Update(context.Background(), tt.appGroupName, tt.appName, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.Status != tt.input.Status {
				t.Errorf("result.Status = %q, want %q", result.Status, tt.input.Status)
			}
		})
	}
}

func TestAppGroupAppService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		appName      string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			appName:      "my-app",
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "app not found",
			appGroupName: "app-group-1",
			appName:      "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "App not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "forbidden",
			appGroupName: "app-group-1",
			appName:      "my-app",
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
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName + "/apps/" + tt.appName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			err := client.AppGroupApps.Delete(context.Background(), tt.appGroupName, tt.appName)
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

func TestAppGroupAppService_List(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		opts         *ListOptions
		mockStatus   int
		mockResponse interface{}
		wantQuery    string
		wantCount    int
		wantErr      bool
	}{
		{
			name:         "success with results",
			appGroupName: "app-group-1",
			opts:         nil,
			mockStatus:   http.StatusOK,
			mockResponse: &AppGroupAppListResponse{
				AppGroupApps: []AppGroupApp{
					{Name: "app-1", Status: "approved"},
					{Name: "app-2", Status: "approved"},
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:         "success with empty results",
			appGroupName: "app-group-1",
			opts:         nil,
			mockStatus:   http.StatusOK,
			mockResponse: &AppGroupAppListResponse{
				AppGroupApps: []AppGroupApp{},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:         "with pagination",
			appGroupName: "app-group-1",
			opts:         &ListOptions{PageSize: 10, PageToken: "token123"},
			mockStatus:   http.StatusOK,
			mockResponse: &AppGroupAppListResponse{
				AppGroupApps: []AppGroupApp{
					{Name: "app-1"},
				},
				NextPageToken: "next-token",
			},
			wantQuery: "pageSize=10&pageToken=token123",
			wantCount: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodGet,
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName + "/apps",
				wantQuery:  tt.wantQuery,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroupApps.List(context.Background(), tt.appGroupName, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(result.AppGroupApps) != tt.wantCount {
				t.Errorf("got %d apps, want %d", len(result.AppGroupApps), tt.wantCount)
			}
		})
	}
}

func TestAppGroupApp_WithCredentials(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/appgroups/app-group-1/apps/my-app",
		response: &AppGroupApp{
			Name:     "my-app",
			AppID:    "app-id-123",
			AppGroup: "app-group-1",
			Status:   "approved",
			Credentials: []AppGroupAppKey{
				{
					ConsumerKey:    "key-1",
					ConsumerSecret: "secret-1",
					Status:         "approved",
					APIProducts: []APIProductRef{
						{APIProduct: "product-1", Status: "approved"},
						{APIProduct: "product-2", Status: "approved"},
					},
				},
				{
					ConsumerKey:    "key-2",
					ConsumerSecret: "secret-2",
					Status:         "revoked",
				},
			},
		},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.AppGroupApps.Get(context.Background(), "app-group-1", "my-app")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if len(result.Credentials) != 2 {
		t.Errorf("Credentials count = %d, want 2", len(result.Credentials))
	}
	if result.Credentials[0].ConsumerKey != "key-1" {
		t.Errorf("Credentials[0].ConsumerKey = %q, want %q", result.Credentials[0].ConsumerKey, "key-1")
	}
	if len(result.Credentials[0].APIProducts) != 2 {
		t.Errorf("Credentials[0].APIProducts count = %d, want 2", len(result.Credentials[0].APIProducts))
	}
}

func TestAppGroupApp_WithAttributes(t *testing.T) {
	input := &AppGroupApp{
		Name:   "app-with-attrs",
		Status: "approved",
		Attributes: []Attribute{
			{Name: "env", Value: "production"},
			{Name: "tier", Value: "premium"},
		},
		Scopes: []string{"read", "write", "admin"},
	}

	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodPost,
		wantPath:   "/organizations/test-org/appgroups/app-group-1/apps",
		response:   input,
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.AppGroupApps.Create(context.Background(), "app-group-1", input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if len(result.Attributes) != 2 {
		t.Errorf("Attributes count = %d, want 2", len(result.Attributes))
	}
	if len(result.Scopes) != 3 {
		t.Errorf("Scopes count = %d, want 3", len(result.Scopes))
	}
}

func TestAppGroupAppListResponse_Pagination(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/appgroups/app-group-1/apps",
		response: &AppGroupAppListResponse{
			AppGroupApps: []AppGroupApp{
				{Name: "app-1"},
				{Name: "app-2"},
			},
			NextPageToken: "page-2-token",
		},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.AppGroupApps.List(context.Background(), "app-group-1", &ListOptions{PageSize: 2})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if result.NextPageToken != "page-2-token" {
		t.Errorf("NextPageToken = %q, want %q", result.NextPageToken, "page-2-token")
	}
}
