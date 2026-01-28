package apigee

import (
	"context"
	"net/http"
	"testing"
)

// Organization-level KVM tests

func TestKeyValueMapService_Create(t *testing.T) {
	tests := []struct {
		name         string
		input        *KeyValueMap
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:       "success",
			input:      &KeyValueMap{Name: "config-kvm", Encrypted: true},
			mockStatus: http.StatusOK,
			mockResponse: &KeyValueMap{
				Name:      "config-kvm",
				Encrypted: true,
			},
			wantErr: false,
		},
		{
			name:         "conflict",
			input:        &KeyValueMap{Name: "existing-kvm"},
			mockStatus:   http.StatusConflict,
			mockResponse: errorResponseBody(409, "KVM already exists", "ALREADY_EXISTS"),
			wantErr:      true,
			errCheck:     IsConflict,
		},
		{
			name:         "unauthorized",
			input:        &KeyValueMap{Name: "config-kvm"},
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
				wantPath:   "/organizations/test-org/keyvaluemaps",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.KeyValueMaps.Create(context.Background(), tt.input)
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

func TestKeyValueMapService_Get(t *testing.T) {
	tests := []struct {
		name         string
		kvmName      string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:       "success",
			kvmName:    "config-kvm",
			mockStatus: http.StatusOK,
			mockResponse: &KeyValueMap{
				Name:      "config-kvm",
				Encrypted: true,
			},
			wantErr: false,
		},
		{
			name:         "not found",
			kvmName:      "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "KVM not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodGet,
				wantPath:   "/organizations/test-org/keyvaluemaps/" + tt.kvmName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.KeyValueMaps.Get(context.Background(), tt.kvmName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("error check failed for error: %v", err)
			}
			if err == nil && result.Name != tt.kvmName {
				t.Errorf("result.Name = %q, want %q", result.Name, tt.kvmName)
			}
		})
	}
}

func TestKeyValueMapService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		kvmName      string
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:       "success",
			kvmName:    "config-kvm",
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "not found",
			kvmName:      "nonexistent",
			mockStatus:   http.StatusNotFound,
			mockResponse: errorResponseBody(404, "KVM not found", "NOT_FOUND"),
			wantErr:      true,
			errCheck:     IsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodDelete,
				wantPath:   "/organizations/test-org/keyvaluemaps/" + tt.kvmName,
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			err := client.KeyValueMaps.Delete(context.Background(), tt.kvmName)
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

func TestKeyValueMapService_List(t *testing.T) {
	tests := []struct {
		name         string
		mockStatus   int
		mockResponse interface{}
		wantNames    []string
		wantErr      bool
	}{
		{
			name:         "success with results",
			mockStatus:   http.StatusOK,
			mockResponse: []string{"kvm-1", "kvm-2", "kvm-3"},
			wantNames:    []string{"kvm-1", "kvm-2", "kvm-3"},
			wantErr:      false,
		},
		{
			name:         "success with empty results",
			mockStatus:   http.StatusOK,
			mockResponse: []string{},
			wantNames:    []string{},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodGet,
				wantPath:   "/organizations/test-org/keyvaluemaps",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.KeyValueMaps.List(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if len(result.KeyValueMapNames) != len(tt.wantNames) {
					t.Errorf("got %d names, want %d", len(result.KeyValueMapNames), len(tt.wantNames))
				}
			}
		})
	}
}

// Organization-level KVM Entry tests

func TestKeyValueMapEntryService_Create(t *testing.T) {
	tests := []struct {
		name         string
		kvmName      string
		input        *KeyValueEntry
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
	}{
		{
			name:       "success",
			kvmName:    "config-kvm",
			input:      &KeyValueEntry{Name: "key1", Value: "value1"},
			mockStatus: http.StatusOK,
			mockResponse: &KeyValueEntry{
				Name:  "key1",
				Value: "value1",
			},
			wantErr: false,
		},
		{
			name:         "conflict",
			kvmName:      "config-kvm",
			input:        &KeyValueEntry{Name: "existing-key", Value: "value"},
			mockStatus:   http.StatusConflict,
			mockResponse: errorResponseBody(409, "Entry already exists", "ALREADY_EXISTS"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodPost,
				wantPath:   "/organizations/test-org/keyvaluemaps/" + tt.kvmName + "/entries",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.KeyValueMapEntries.Create(context.Background(), tt.kvmName, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && result.Name != tt.input.Name {
				t.Errorf("result.Name = %q, want %q", result.Name, tt.input.Name)
			}
		})
	}
}

func TestKeyValueMapEntryService_Get(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/keyvaluemaps/config-kvm/entries/key1",
		response:   &KeyValueEntry{Name: "key1", Value: "value1"},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.KeyValueMapEntries.Get(context.Background(), "config-kvm", "key1")
	if err != nil {
		t.Errorf("Get() error = %v", err)
		return
	}
	if result.Name != "key1" {
		t.Errorf("result.Name = %q, want %q", result.Name, "key1")
	}
}

func TestKeyValueMapEntryService_Update(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodPut,
		wantPath:   "/organizations/test-org/keyvaluemaps/config-kvm/entries/key1",
		response:   &KeyValueEntry{Name: "key1", Value: "updated-value"},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	input := &KeyValueEntry{Name: "key1", Value: "updated-value"}
	result, err := client.KeyValueMapEntries.Update(context.Background(), "config-kvm", "key1", input)
	if err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}
	if result.Value != "updated-value" {
		t.Errorf("result.Value = %q, want %q", result.Value, "updated-value")
	}
}

func TestKeyValueMapEntryService_Delete(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodDelete,
		wantPath:   "/organizations/test-org/keyvaluemaps/config-kvm/entries/key1",
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	err := client.KeyValueMapEntries.Delete(context.Background(), "config-kvm", "key1")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}
}

func TestKeyValueMapEntryService_List(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/keyvaluemaps/config-kvm/entries",
		response: &KeyValueEntryListResponse{
			KeyValueEntries: []KeyValueEntry{
				{Name: "key1", Value: "value1"},
				{Name: "key2", Value: "value2"},
			},
		},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.KeyValueMapEntries.List(context.Background(), "config-kvm")
	if err != nil {
		t.Errorf("List() error = %v", err)
		return
	}
	if len(result.KeyValueEntries) != 2 {
		t.Errorf("got %d entries, want 2", len(result.KeyValueEntries))
	}
}

// Environment-level KVM tests

func TestEnvKeyValueMapService_Create(t *testing.T) {
	tests := []struct {
		name         string
		envName      string
		input        *KeyValueMap
		mockStatus   int
		mockResponse interface{}
		wantErr      bool
		errCheck     func(error) bool
	}{
		{
			name:       "success",
			envName:    "test-env",
			input:      &KeyValueMap{Name: "env-config", Encrypted: true},
			mockStatus: http.StatusOK,
			mockResponse: &KeyValueMap{
				Name:      "env-config",
				Encrypted: true,
			},
			wantErr: false,
		},
		{
			name:         "conflict",
			envName:      "test-env",
			input:        &KeyValueMap{Name: "existing"},
			mockStatus:   http.StatusConflict,
			mockResponse: errorResponseBody(409, "KVM already exists", "ALREADY_EXISTS"),
			wantErr:      true,
			errCheck:     IsConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &requestValidator{
				t:          t,
				wantMethod: http.MethodPost,
				wantPath:   "/organizations/test-org/environments/" + tt.envName + "/keyvaluemaps",
				response:   tt.mockResponse,
				statusCode: tt.mockStatus,
			}
			client, _ := setupTestClient(t, handler)

			result, err := client.EnvKeyValueMaps.Create(context.Background(), tt.envName, tt.input)
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

func TestEnvKeyValueMapService_Get(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/environments/test-env/keyvaluemaps/env-config",
		response:   &KeyValueMap{Name: "env-config", Encrypted: true},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.EnvKeyValueMaps.Get(context.Background(), "test-env", "env-config")
	if err != nil {
		t.Errorf("Get() error = %v", err)
		return
	}
	if result.Name != "env-config" {
		t.Errorf("result.Name = %q, want %q", result.Name, "env-config")
	}
}

func TestEnvKeyValueMapService_Delete(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodDelete,
		wantPath:   "/organizations/test-org/environments/test-env/keyvaluemaps/env-config",
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	err := client.EnvKeyValueMaps.Delete(context.Background(), "test-env", "env-config")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}
}

func TestEnvKeyValueMapService_List(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/environments/test-env/keyvaluemaps",
		response:   []string{"kvm-1", "kvm-2"},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.EnvKeyValueMaps.List(context.Background(), "test-env")
	if err != nil {
		t.Errorf("List() error = %v", err)
		return
	}
	if len(result.KeyValueMapNames) != 2 {
		t.Errorf("got %d names, want 2", len(result.KeyValueMapNames))
	}
}

// Environment-level KVM Entry tests

func TestEnvKeyValueMapEntryService_Create(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodPost,
		wantPath:   "/organizations/test-org/environments/test-env/keyvaluemaps/env-config/entries",
		response:   &KeyValueEntry{Name: "key1", Value: "value1"},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	input := &KeyValueEntry{Name: "key1", Value: "value1"}
	result, err := client.EnvKeyValueMapEntries.Create(context.Background(), "test-env", "env-config", input)
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}
	if result.Name != "key1" {
		t.Errorf("result.Name = %q, want %q", result.Name, "key1")
	}
}

func TestEnvKeyValueMapEntryService_Get(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/environments/test-env/keyvaluemaps/env-config/entries/key1",
		response:   &KeyValueEntry{Name: "key1", Value: "value1"},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.EnvKeyValueMapEntries.Get(context.Background(), "test-env", "env-config", "key1")
	if err != nil {
		t.Errorf("Get() error = %v", err)
		return
	}
	if result.Name != "key1" {
		t.Errorf("result.Name = %q, want %q", result.Name, "key1")
	}
}

func TestEnvKeyValueMapEntryService_Update(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodPut,
		wantPath:   "/organizations/test-org/environments/test-env/keyvaluemaps/env-config/entries/key1",
		response:   &KeyValueEntry{Name: "key1", Value: "updated"},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	input := &KeyValueEntry{Name: "key1", Value: "updated"}
	result, err := client.EnvKeyValueMapEntries.Update(context.Background(), "test-env", "env-config", "key1", input)
	if err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}
	if result.Value != "updated" {
		t.Errorf("result.Value = %q, want %q", result.Value, "updated")
	}
}

func TestEnvKeyValueMapEntryService_Delete(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodDelete,
		wantPath:   "/organizations/test-org/environments/test-env/keyvaluemaps/env-config/entries/key1",
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	err := client.EnvKeyValueMapEntries.Delete(context.Background(), "test-env", "env-config", "key1")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}
}

func TestEnvKeyValueMapEntryService_List(t *testing.T) {
	handler := &requestValidator{
		t:          t,
		wantMethod: http.MethodGet,
		wantPath:   "/organizations/test-org/environments/test-env/keyvaluemaps/env-config/entries",
		response: &KeyValueEntryListResponse{
			KeyValueEntries: []KeyValueEntry{
				{Name: "key1", Value: "value1"},
				{Name: "key2", Value: "value2"},
			},
		},
		statusCode: http.StatusOK,
	}
	client, _ := setupTestClient(t, handler)

	result, err := client.EnvKeyValueMapEntries.List(context.Background(), "test-env", "env-config")
	if err != nil {
		t.Errorf("List() error = %v", err)
		return
	}
	if len(result.KeyValueEntries) != 2 {
		t.Errorf("got %d entries, want 2", len(result.KeyValueEntries))
	}
}
