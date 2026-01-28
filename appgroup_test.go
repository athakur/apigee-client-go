package apigee

import (
	"context"
	"net/http"
	"testing"
)

func TestAppGroupService_Create(t *testing.T) {
	tests := []struct {
		name         string
		input        *AppGroup
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name: "success",
			input: &AppGroup{
				Name:        "app-group-1",
				DisplayName: "App Group One",
				ChannelID:   "channel-1",
				ChannelURI:  "https://example.com/channel",
				Status:      "active",
			},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroup{
				Name:        "app-group-1",
				AppGroupID:  "generated-id",
				DisplayName: "App Group One",
				ChannelID:   "channel-1",
				ChannelURI:  "https://example.com/channel",
				Status:      "active",
			},
			wantErr: false,
		},
		{
			name:         "conflict",
			input:        &AppGroup{Name: "existing"},
			mockStatus:   http.StatusConflict,
			mockResponse: errorResponseBody(409, "App group already exists", "ALREADY_EXISTS"),
			wantErr:      true,
			errCheck:     IsConflict,
		},
		{
			name:         "unauthorized",
			input:        &AppGroup{Name: "app-group-1"},
			mockStatus:   http.StatusUnauthorized,
			mockResponse: errorResponseBody(401, "Invalid credentials", "UNAUTHENTICATED"),
			wantErr:      true,
			errCheck:     IsUnauthorized,
		},
		{
			name:         "forbidden",
			input:        &AppGroup{Name: "app-group-1"},
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
				wantMethod: http.MethodPost,
				wantPath:   "/organizations/test-org/appgroups",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroups.Create(context.Background(), tt.input)
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

func TestAppGroupService_Get(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			mockStatus:   http.StatusOK,
			mockResponse: &AppGroup{
				Name:        "app-group-1",
				AppGroupID:  "id-123",
				DisplayName: "App Group One",
				Status:      "active",
			},
			wantErr: false,
		},
		{
			name:         "not found",
			appGroupName: "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "App group not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "unauthorized",
			appGroupName: "app-group-1",
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
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroups.Get(context.Background(), tt.appGroupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.Name != tt.appGroupName {
				t.Errorf("result.Name = %q, want %q", result.Name, tt.appGroupName)
			}
		})
	}
}

func TestAppGroupService_Update(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		input        *AppGroup
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			input: &AppGroup{
				Name:        "app-group-1",
				DisplayName: "Updated App Group",
				Status:      "inactive",
			},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroup{
				Name:        "app-group-1",
				DisplayName: "Updated App Group",
				Status:      "inactive",
			},
			wantErr: false,
		},
		{
			name:         "not found",
			appGroupName: "nonexistent",
			input: &AppGroup{
				Name:        "nonexistent",
				DisplayName: "Test",
			},
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "App group not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodPut,
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroups.Update(context.Background(), tt.appGroupName, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.DisplayName != tt.input.DisplayName {
				t.Errorf("result.DisplayName = %q, want %q", result.DisplayName, tt.input.DisplayName)
			}
		})
	}
}

func TestAppGroupService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		appGroupName string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success",
			appGroupName: "app-group-1",
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "not found",
			appGroupName: "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "App group not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "forbidden",
			appGroupName: "app-group-1",
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
				wantPath:   "/organizations/test-org/appgroups/" + tt.appGroupName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			err := client.AppGroups.Delete(context.Background(), tt.appGroupName)
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

func TestAppGroupService_List(t *testing.T) {
	tests := []struct {
		name         string
		opts         *ListOptions
		mockStatus   int
		mockResponse interface{}
		wantQuery    string
		wantCount    int
		wantErr      bool
	}{
		{
			name:       "success with results",
			opts:       nil,
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupListResponse{
				AppGroups: []AppGroup{
					{Name: "app-group-1", DisplayName: "App Group 1"},
					{Name: "app-group-2", DisplayName: "App Group 2"},
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:       "success with empty results",
			opts:       nil,
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupListResponse{
				AppGroups: []AppGroup{},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:       "with pagination",
			opts:       &ListOptions{PageSize: 10, PageToken: "token123"},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupListResponse{
				AppGroups: []AppGroup{
					{Name: "app-group-1"},
				},
				NextPageToken: "next-token",
			},
			wantQuery: "pageSize=10&pageToken=token123",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:       "with filter",
			opts:       &ListOptions{Filter: "status=active"},
			mockStatus: http.StatusOK,
			mockResponse: &AppGroupListResponse{
				AppGroups: []AppGroup{
					{Name: "app-group-1", Status: "active"},
				},
			},
			wantQuery: "filter=status%3Dactive",
			wantCount: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodGet,
				wantPath:   "/organizations/test-org/appgroups",
				wantQuery:  tt.wantQuery,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.AppGroups.List(context.Background(), tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(result.AppGroups) != tt.wantCount {
				t.Errorf("got %d app groups, want %d", len(result.AppGroups), tt.wantCount)
			}
		})
	}
}

func TestAppGroup_WithAttributes(t *testing.T) {
	input := &AppGroup{
		Name:        "app-group-with-attrs",
		DisplayName: "App Group With Attributes",
		Status:      "active",
		Attributes: []Attribute{
			{Name: "env", Value: "production"},
			{Name: "tier", Value: "premium"},
			{Name: "region", Value: "us-east"},
		},
	}

	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodPost,
		wantPath:   "/organizations/test-org/appgroups",
		response:   input,
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.AppGroups.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if len(result.Attributes) != 3 {
		t.Errorf("Attributes count = %d, want 3", len(result.Attributes))
	}
	if result.Attributes[0].Name != "env" {
		t.Errorf("Attributes[0].Name = %q, want %q", result.Attributes[0].Name, "env")
	}
}

func TestAppGroupListResponse_Pagination(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/appgroups",
		response: &AppGroupListResponse{
			AppGroups: []AppGroup{
				{Name: "app-group-1"},
				{Name: "app-group-2"},
			},
			NextPageToken: "page-2-token",
		},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.AppGroups.List(context.Background(), &ListOptions{PageSize: 2})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if result.NextPageToken != "page-2-token" {
		t.Errorf("NextPageToken = %q, want %q", result.NextPageToken, "page-2-token")
	}
}
