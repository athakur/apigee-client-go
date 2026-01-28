package apigee

import (
	"context"
	"net/http"
	"testing"
)

func TestAPIProductService_Create(t *testing.T) {
	tests := []struct {
		name         string
		input        *APIProduct
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name: "success",
			input: &APIProduct{
				Name:         "product-1",
				DisplayName:  "Product One",
				Description:  "Test product",
				ApprovalType: "auto",
				Environments: []string{"test"},
				Proxies:      []string{"proxy-1"},
			},
			mockStatus: http.StatusOK,
			mockResponse: &APIProduct{
				Name:         "product-1",
				DisplayName:  "Product One",
				Description:  "Test product",
				ApprovalType: "auto",
				Environments: []string{"test"},
				Proxies:      []string{"proxy-1"},
			},
			wantErr: false,
		},
		{
			name:         "conflict",
			input:        &APIProduct{Name: "existing"},
			mockStatus:   http.StatusConflict,
			mockResponse: errorResponseBody(409, "Product already exists", "ALREADY_EXISTS"),
			wantErr:      true,
			errCheck:     IsConflict,
		},
		{
			name:         "unauthorized",
			input:        &APIProduct{Name: "product-1"},
			mockStatus:   http.StatusUnauthorized,
			mockResponse: errorResponseBody(401, "Invalid credentials", "UNAUTHENTICATED"),
			wantErr:      true,
			errCheck:     IsUnauthorized,
		},
		{
			name:         "forbidden",
			input:        &APIProduct{Name: "product-1"},
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
				wantPath:   "/organizations/test-org/apiproducts",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.APIProducts.Create(context.Background(), tt.input)
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

func TestAPIProductService_Get(t *testing.T) {
	tests := []struct {
		name         string
		productName  string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:        "success",
			productName: "product-1",
			mockStatus:  http.StatusOK,
			mockResponse: &APIProduct{
				Name:         "product-1",
				DisplayName:  "Product One",
				ApprovalType: "auto",
			},
			wantErr: false,
		},
		{
			name:         "not found",
			productName:  "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Product not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "unauthorized",
			productName:  "product-1",
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
				wantPath:   "/organizations/test-org/apiproducts/" + tt.productName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.APIProducts.Get(context.Background(), tt.productName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.Name != tt.productName {
				t.Errorf("result.Name = %q, want %q", result.Name, tt.productName)
			}
		})
	}
}

func TestAPIProductService_Update(t *testing.T) {
	tests := []struct {
		name         string
		productName  string
		input        *APIProduct
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:        "success",
			productName: "product-1",
			input: &APIProduct{
				Name:         "product-1",
				DisplayName:  "Updated Product",
				ApprovalType: "manual",
			},
			mockStatus: http.StatusOK,
			mockResponse: &APIProduct{
				Name:         "product-1",
				DisplayName:  "Updated Product",
				ApprovalType: "manual",
			},
			wantErr: false,
		},
		{
			name:        "not found",
			productName: "nonexistent",
			input: &APIProduct{
				Name:        "nonexistent",
				DisplayName: "Test",
			},
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Product not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodPut,
				wantPath:   "/organizations/test-org/apiproducts/" + tt.productName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.APIProducts.Update(context.Background(), tt.productName, tt.input)
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

func TestAPIProductService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		productName  string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:        "success",
			productName: "product-1",
			mockStatus:  http.StatusOK,
			wantErr:     false,
		},
		{
			name:         "not found",
			productName:  "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Product not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "forbidden",
			productName:  "product-1",
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
				wantPath:   "/organizations/test-org/apiproducts/" + tt.productName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			err := client.APIProducts.Delete(context.Background(), tt.productName)
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

func TestAPIProductService_List(t *testing.T) {
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
			mockResponse: &APIProductListResponse{
				APIProducts: []APIProduct{
					{Name: "product-1", DisplayName: "Product 1"},
					{Name: "product-2", DisplayName: "Product 2"},
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:       "success with empty results",
			opts:       nil,
			mockStatus: http.StatusOK,
			mockResponse: &APIProductListResponse{
				APIProducts: []APIProduct{},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:       "with pagination",
			opts:       &ListOptions{PageSize: 10, PageToken: "token123"},
			mockStatus: http.StatusOK,
			mockResponse: &APIProductListResponse{
				APIProducts: []APIProduct{
					{Name: "product-1"},
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
			mockResponse: &APIProductListResponse{
				APIProducts: []APIProduct{
					{Name: "product-1"},
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
				wantPath:   "/organizations/test-org/apiproducts",
				wantQuery:  tt.wantQuery,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.APIProducts.List(context.Background(), tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(result.APIProducts) != tt.wantCount {
				t.Errorf("got %d products, want %d", len(result.APIProducts), tt.wantCount)
			}
		})
	}
}

func TestAPIProduct_WithAllFields(t *testing.T) {
	input := &APIProduct{
		Name:          "complete-product",
		DisplayName:   "Complete Product",
		Description:   "A product with all fields",
		ApprovalType:  "auto",
		Environments:  []string{"test", "prod"},
		Proxies:       []string{"proxy-1", "proxy-2"},
		APIResources:  []string{"/users/**", "/orders/**"},
		Scopes:        []string{"read", "write"},
		Quota:         "100",
		QuotaInterval: "1",
		QuotaTimeUnit: "hour",
		Attributes: []Attribute{
			{Name: "env", Value: "production"},
			{Name: "tier", Value: "premium"},
		},
	}

	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodPost,
		wantPath:   "/organizations/test-org/apiproducts",
		response:   input,
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.APIProducts.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if result.Name != input.Name {
		t.Errorf("Name = %q, want %q", result.Name, input.Name)
	}
	if len(result.Environments) != 2 {
		t.Errorf("Environments count = %d, want 2", len(result.Environments))
	}
	if len(result.Proxies) != 2 {
		t.Errorf("Proxies count = %d, want 2", len(result.Proxies))
	}
	if len(result.Scopes) != 2 {
		t.Errorf("Scopes count = %d, want 2", len(result.Scopes))
	}
	if len(result.Attributes) != 2 {
		t.Errorf("Attributes count = %d, want 2", len(result.Attributes))
	}
}

func TestAPIProductListResponse_Pagination(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/apiproducts",
		response: &APIProductListResponse{
			APIProducts: []APIProduct{
				{Name: "product-1"},
				{Name: "product-2"},
			},
			NextPageToken: "page-2-token",
		},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.APIProducts.List(context.Background(), &ListOptions{PageSize: 2})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if result.NextPageToken != "page-2-token" {
		t.Errorf("NextPageToken = %q, want %q", result.NextPageToken, "page-2-token")
	}
}
