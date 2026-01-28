package apigee

import (
	"context"
	"net/http"
	"testing"
)

func TestTargetServer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ts      *TargetServer
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid target server",
			ts: &TargetServer{
				Name: "backend-1",
				Host: "api.example.com",
				Port: 443,
			},
			wantErr: false,
		},
		{
			name: "valid with all fields",
			ts: &TargetServer{
				Name:        "backend-1",
				Host:        "api.example.com",
				Port:        8080,
				Protocol:    "HTTP",
				IsEnabled:   true,
				Description: "Backend server",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			ts: &TargetServer{
				Host: "api.example.com",
				Port: 443,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing host",
			ts: &TargetServer{
				Name: "backend-1",
				Port: 443,
			},
			wantErr: true,
			errMsg:  "host is required",
		},
		{
			name: "port zero",
			ts: &TargetServer{
				Name: "backend-1",
				Host: "api.example.com",
				Port: 0,
			},
			wantErr: true,
			errMsg:  "port must be between 1 and 65535",
		},
		{
			name: "port negative",
			ts: &TargetServer{
				Name: "backend-1",
				Host: "api.example.com",
				Port: -1,
			},
			wantErr: true,
			errMsg:  "port must be between 1 and 65535",
		},
		{
			name: "port too high",
			ts: &TargetServer{
				Name: "backend-1",
				Host: "api.example.com",
				Port: 65536,
			},
			wantErr: true,
			errMsg:  "port must be between 1 and 65535",
		},
		{
			name: "min valid port",
			ts: &TargetServer{
				Name: "backend-1",
				Host: "api.example.com",
				Port: 1,
			},
			wantErr: false,
		},
		{
			name: "max valid port",
			ts: &TargetServer{
				Name: "backend-1",
				Host: "api.example.com",
				Port: 65535,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if got := err.Error(); got != "apigee: target server "+tt.errMsg {
					t.Errorf("error = %q, want containing %q", got, tt.errMsg)
				}
			}
		})
	}
}

func TestTargetServerService_Create(t *testing.T) {
	tests := []struct {
		name         string
		envName      string
		input        *TargetServer
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:    "success",
			envName: "test-env",
			input: &TargetServer{
				Name:      "backend-1",
				Host:      "api.example.com",
				Port:      443,
				IsEnabled: true,
			},
			mockStatus: http.StatusOK,
			mockResponse: &TargetServer{
				Name:      "backend-1",
				Host:      "api.example.com",
				Port:      443,
				IsEnabled: true,
			},
			wantErr: false,
		},
		{
			name:    "validation error - missing name",
			envName: "test-env",
			input: &TargetServer{
				Host: "api.example.com",
				Port: 443,
			},
			wantErr: true,
		},
		{
			name:    "validation error - invalid port",
			envName: "test-env",
			input: &TargetServer{
				Name: "backend-1",
				Host: "api.example.com",
				Port: 0,
			},
			wantErr: true,
		},
		{
			name:    "conflict - already exists",
			envName: "test-env",
			input: &TargetServer{
				Name: "existing",
				Host: "api.example.com",
				Port: 443,
			},
			mockStatus:   http.StatusConflict,
			mockResponse: errorResponseBody(409, "Target server already exists", "ALREADY_EXISTS"),
			wantErr:      true,
			errCheck:     IsConflict,
		},
		{
			name:    "unauthorized",
			envName: "test-env",
			input: &TargetServer{
				Name: "backend-1",
				Host: "api.example.com",
				Port: 443,
			},
			mockStatus:   http.StatusUnauthorized,
			mockResponse: errorResponseBody(401, "Invalid credentials", "UNAUTHENTICATED"),
			wantErr:      true,
			errCheck:     IsUnauthorized,
		},
		{
			name:    "forbidden",
			envName: "test-env",
			input: &TargetServer{
				Name: "backend-1",
				Host: "api.example.com",
				Port: 443,
			},
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
				wantPath:   "/organizations/test-org/environments/" + tt.envName + "/targetservers",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.TargetServers.Create(context.Background(), tt.envName, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil {
				if result.Name != tt.input.Name {
					t.Errorf("result.Name = %q, want %q", result.Name, tt.input.Name)
				}
			}
		})
	}
}

func TestTargetServerService_Get(t *testing.T) {
	tests := []struct {
		name         string
		envName      string
		tsName       string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:       "success",
			envName:    "test-env",
			tsName:     "backend-1",
			mockStatus: http.StatusOK,
			mockResponse: &TargetServer{
				Name:      "backend-1",
				Host:      "api.example.com",
				Port:      443,
				IsEnabled: true,
			},
			wantErr: false,
		},
		{
			name:         "not found",
			envName:      "test-env",
			tsName:       "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Target server not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "unauthorized",
			envName:      "test-env",
			tsName:       "backend-1",
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
				wantPath:   "/organizations/test-org/environments/" + tt.envName + "/targetservers/" + tt.tsName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.TargetServers.Get(context.Background(), tt.envName, tt.tsName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.Name != tt.tsName {
				t.Errorf("result.Name = %q, want %q", result.Name, tt.tsName)
			}
		})
	}
}

func TestTargetServerService_Update(t *testing.T) {
	tests := []struct {
		name         string
		envName      string
		tsName       string
		input        *TargetServer
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:    "success",
			envName: "test-env",
			tsName:  "backend-1",
			input: &TargetServer{
				Name:      "backend-1",
				Host:      "updated.example.com",
				Port:      8443,
				IsEnabled: false,
			},
			mockStatus: http.StatusOK,
			mockResponse: &TargetServer{
				Name:      "backend-1",
				Host:      "updated.example.com",
				Port:      8443,
				IsEnabled: false,
			},
			wantErr: false,
		},
		{
			name:    "validation error",
			envName: "test-env",
			tsName:  "backend-1",
			input: &TargetServer{
				Name: "backend-1",
				Host: "",
				Port: 443,
			},
			wantErr: true,
		},
		{
			name:    "not found",
			envName: "test-env",
			tsName:  "nonexistent",
			input: &TargetServer{
				Name: "nonexistent",
				Host: "api.example.com",
				Port: 443,
			},
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Target server not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodPut,
				wantPath:   "/organizations/test-org/environments/" + tt.envName + "/targetservers/" + tt.tsName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.TargetServers.Update(context.Background(), tt.envName, tt.tsName, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.Host != tt.input.Host {
				t.Errorf("result.Host = %q, want %q", result.Host, tt.input.Host)
			}
		})
	}
}

func TestTargetServerService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		envName      string
		tsName       string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:       "success",
			envName:    "test-env",
			tsName:     "backend-1",
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "not found",
			envName:      "test-env",
			tsName:       "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "Target server not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
		{
			name:         "forbidden",
			envName:      "test-env",
			tsName:       "backend-1",
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
				wantPath:   "/organizations/test-org/environments/" + tt.envName + "/targetservers/" + tt.tsName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			err := client.TargetServers.Delete(context.Background(), tt.envName, tt.tsName)
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

func TestTargetServerService_List(t *testing.T) {
	tests := []struct {
		name         string
		envName      string
		mockStatus   int
		mockResponse interface{}
		wantNames    []string
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:         "success with results",
			envName:      "test-env",
			mockStatus:   http.StatusOK,
			mockResponse: []string{"backend-1", "backend-2", "backend-3"},
			wantNames:    []string{"backend-1", "backend-2", "backend-3"},
			wantErr:      false,
		},
		{
			name:         "success with empty results",
			envName:      "test-env",
			mockStatus:   http.StatusOK,
			mockResponse: []string{},
			wantNames:    []string{},
			wantErr:      false,
		},
		{
			name:         "unauthorized",
			envName:      "test-env",
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
				wantPath:   "/organizations/test-org/environments/" + tt.envName + "/targetservers",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.TargetServers.List(context.Background(), tt.envName)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil {
				if len(result.TargetServerNames) != len(tt.wantNames) {
					t.Errorf("got %d names, want %d", len(result.TargetServerNames), len(tt.wantNames))
				}
				for i, name := range result.TargetServerNames {
					if name != tt.wantNames[i] {
						t.Errorf("name[%d] = %q, want %q", i, name, tt.wantNames[i])
					}
				}
			}
		})
	}
}

func TestTargetServerService_Create_withSSL(t *testing.T) {
	input := &TargetServer{
		Name:      "secure-backend",
		Host:      "api.example.com",
		Port:      443,
		IsEnabled: true,
		SSLInfo: &SSLInfo{
			Enabled:           true,
			ClientAuthEnabled: true,
			KeyStore:          "my-keystore",
			TrustStore:        "my-truststore",
			KeyAlias:          "my-key",
			Protocols:         []string{"TLSv1.2", "TLSv1.3"},
			CommonName: &CommonName{
				Value:         "*.example.com",
				WildcardMatch: true,
			},
		},
	}

	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodPost,
		wantPath:   "/organizations/test-org/environments/test-env/targetservers",
		response:   input,
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.TargetServers.Create(context.Background(), "test-env", input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if result.SSLInfo == nil {
		t.Fatal("SSLInfo should not be nil")
	}
	if !result.SSLInfo.Enabled {
		t.Error("SSLInfo.Enabled = false, want true")
	}
	if !result.SSLInfo.ClientAuthEnabled {
		t.Error("SSLInfo.ClientAuthEnabled = false, want true")
	}
	if result.SSLInfo.CommonName == nil {
		t.Fatal("SSLInfo.CommonName should not be nil")
	}
	if !result.SSLInfo.CommonName.WildcardMatch {
		t.Error("SSLInfo.CommonName.WildcardMatch = false, want true")
	}
}
