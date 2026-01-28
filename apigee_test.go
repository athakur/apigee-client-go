package apigee

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name         string
		organization string
		opts         []ClientOption
		wantErr      bool
		errContains  string
		validate     func(*testing.T, *Client)
	}{
		{
			name:         "valid organization with token source",
			organization: "my-org",
			opts:         []ClientOption{WithTokenSource(&mockTokenSource{})},
			wantErr:      false,
			validate: func(t *testing.T, c *Client) {
				if c.Organization != "my-org" {
					t.Errorf("Organization = %q, want %q", c.Organization, "my-org")
				}
				if c.BaseURL != DefaultBaseURL {
					t.Errorf("BaseURL = %q, want %q", c.BaseURL, DefaultBaseURL)
				}
			},
		},
		{
			name:         "empty organization",
			organization: "",
			opts:         []ClientOption{WithTokenSource(&mockTokenSource{})},
			wantErr:      true,
			errContains:  "organization is required",
		},
		{
			name:         "custom base URL",
			organization: "my-org",
			opts: []ClientOption{
				WithTokenSource(&mockTokenSource{}),
				WithBaseURL("https://custom.example.com/v1"),
			},
			wantErr: false,
			validate: func(t *testing.T, c *Client) {
				if c.BaseURL != "https://custom.example.com/v1" {
					t.Errorf("BaseURL = %q, want %q", c.BaseURL, "https://custom.example.com/v1")
				}
			},
		},
		{
			name:         "custom HTTP client",
			organization: "my-org",
			opts: []ClientOption{
				WithTokenSource(&mockTokenSource{}),
				WithHTTPClient(&http.Client{Timeout: 30}),
			},
			wantErr: false,
			validate: func(t *testing.T, c *Client) {
				if c.httpClient == nil {
					t.Error("httpClient should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(context.Background(), tt.organization, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %q, want containing %q", err.Error(), tt.errContains)
				}
				return
			}
			if tt.validate != nil {
				tt.validate(t, client)
			}
		})
	}
}

func TestNewClient_servicesInitialized(t *testing.T) {
	client, err := NewClient(context.Background(), "test-org",
		WithTokenSource(&mockTokenSource{}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// Verify all services are initialized
	if client.AppGroups == nil {
		t.Error("AppGroups service should be initialized")
	}
	if client.AppGroupApps == nil {
		t.Error("AppGroupApps service should be initialized")
	}
	if client.AppGroupAppKeys == nil {
		t.Error("AppGroupAppKeys service should be initialized")
	}
	if client.APIProducts == nil {
		t.Error("APIProducts service should be initialized")
	}
	if client.KeyValueMaps == nil {
		t.Error("KeyValueMaps service should be initialized")
	}
	if client.KeyValueMapEntries == nil {
		t.Error("KeyValueMapEntries service should be initialized")
	}
	if client.EnvKeyValueMaps == nil {
		t.Error("EnvKeyValueMaps service should be initialized")
	}
	if client.EnvKeyValueMapEntries == nil {
		t.Error("EnvKeyValueMapEntries service should be initialized")
	}
	if client.TargetServers == nil {
		t.Error("TargetServers service should be initialized")
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 60}
	client, err := NewClient(context.Background(), "test-org",
		WithTokenSource(&mockTokenSource{}),
		WithHTTPClient(customClient),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client.httpClient != customClient {
		t.Error("httpClient was not set to custom client")
	}
}

func TestWithBaseURL(t *testing.T) {
	customURL := "https://my-apigee.example.com/v1"
	client, err := NewClient(context.Background(), "test-org",
		WithTokenSource(&mockTokenSource{}),
		WithBaseURL(customURL),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client.BaseURL != customURL {
		t.Errorf("BaseURL = %q, want %q", client.BaseURL, customURL)
	}
}

func TestWithTokenSource(t *testing.T) {
	customToken := &oauth2.Token{AccessToken: "custom-token"}
	ts := &mockTokenSource{token: customToken}

	client, err := NewClient(context.Background(), "test-org",
		WithTokenSource(ts),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client.tokenSource != ts {
		t.Error("tokenSource was not set to custom source")
	}

	// Verify the token is returned correctly
	token, err := client.tokenSource.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if token.AccessToken != "custom-token" {
		t.Errorf("AccessToken = %q, want %q", token.AccessToken, "custom-token")
	}
}

func TestClientConstants(t *testing.T) {
	if DefaultBaseURL != "https://apigee.googleapis.com/v1" {
		t.Errorf("DefaultBaseURL = %q, want %q", DefaultBaseURL, "https://apigee.googleapis.com/v1")
	}
	if CloudPlatformScope != "https://www.googleapis.com/auth/cloud-platform" {
		t.Errorf("CloudPlatformScope = %q, want %q", CloudPlatformScope, "https://www.googleapis.com/auth/cloud-platform")
	}
}

func TestListOptions(t *testing.T) {
	opts := &ListOptions{
		PageSize:  50,
		PageToken: "next-page-token",
		Filter:    "status=active",
	}

	if opts.PageSize != 50 {
		t.Errorf("PageSize = %d, want 50", opts.PageSize)
	}
	if opts.PageToken != "next-page-token" {
		t.Errorf("PageToken = %q, want %q", opts.PageToken, "next-page-token")
	}
	if opts.Filter != "status=active" {
		t.Errorf("Filter = %q, want %q", opts.Filter, "status=active")
	}
}

func TestAttribute(t *testing.T) {
	attr := Attribute{
		Name:  "env",
		Value: "production",
	}

	if attr.Name != "env" {
		t.Errorf("Name = %q, want %q", attr.Name, "env")
	}
	if attr.Value != "production" {
		t.Errorf("Value = %q, want %q", attr.Value, "production")
	}
}
