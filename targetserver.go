package apigee

import (
	"context"
	"fmt"
)

// TargetServer represents an Apigee Target Server.
// Target servers define backend server endpoints for load balancing and failover.
type TargetServer struct {
	// Name is the target server identifier.
	Name string `json:"name,omitempty"`
	// Host is the backend hostname or IP address.
	Host string `json:"host,omitempty"`
	// Port is the connection port (1-65535).
	Port int `json:"port,omitempty"`
	// Protocol is the connection protocol (HTTP, HTTP2, GRPC, GRPC_TARGET, EXTERNAL_CALLOUT).
	Protocol string `json:"protocol,omitempty"`
	// IsEnabled indicates if the target server is enabled.
	IsEnabled bool `json:"isEnabled"`
	// Description is a human-readable description.
	Description string `json:"description,omitempty"`
	// SSLInfo contains TLS/SSL configuration.
	SSLInfo *SSLInfo `json:"sSLInfo,omitempty"`
}

// Validate validates the target server fields.
func (ts *TargetServer) Validate() error {
	if ts.Name == "" {
		return fmt.Errorf("apigee: target server name is required")
	}
	if ts.Host == "" {
		return fmt.Errorf("apigee: target server host is required")
	}
	if ts.Port < 1 || ts.Port > 65535 {
		return fmt.Errorf("apigee: target server port must be between 1 and 65535")
	}
	return nil
}

// SSLInfo represents TLS/SSL configuration for a target server.
type SSLInfo struct {
	// Enabled indicates if TLS is enabled.
	Enabled bool `json:"enabled"`
	// ClientAuthEnabled indicates if two-way TLS (mTLS) is enabled.
	ClientAuthEnabled bool `json:"clientAuthEnabled"`
	// IgnoreValidationErrors indicates if certificate validation errors should be ignored.
	IgnoreValidationErrors bool `json:"ignoreValidationErrors"`
	// KeyAlias is the private key/certificate alias.
	KeyAlias string `json:"keyAlias,omitempty"`
	// KeyStore is the keystore resource ID.
	KeyStore string `json:"keyStore,omitempty"`
	// TrustStore is the truststore resource ID.
	TrustStore string `json:"trustStore,omitempty"`
	// Protocols is the list of TLS protocol versions.
	Protocols []string `json:"protocols,omitempty"`
	// Ciphers is the list of cipher suites.
	Ciphers []string `json:"ciphers,omitempty"`
	// CommonName contains certificate common name configuration.
	CommonName *CommonName `json:"commonName,omitempty"`
}

// CommonName represents the certificate common name configuration.
type CommonName struct {
	// Value is the common name value.
	Value string `json:"value,omitempty"`
	// WildcardMatch indicates if wildcard matching is allowed.
	WildcardMatch bool `json:"wildcardMatch"`
}

// TargetServerListResponse is the response for listing target servers.
// Note: The Apigee API returns a simple array of target server names.
type TargetServerListResponse struct {
	// TargetServerNames is the list of target server names.
	TargetServerNames []string
}

// TargetServerService handles operations on environment-level Target Servers.
type TargetServerService struct {
	client *Client
}

// Create creates a new target server in the specified environment.
func (s *TargetServerService) Create(ctx context.Context, envName string, ts *TargetServer) (*TargetServer, error) {
	if err := ts.Validate(); err != nil {
		return nil, err
	}
	return doCreate[TargetServer](ctx, s.client, s.client.envPath(envName, "targetservers"), ts)
}

// Get retrieves a target server by name from the specified environment.
func (s *TargetServerService) Get(ctx context.Context, envName, name string) (*TargetServer, error) {
	return doGet[TargetServer](ctx, s.client, s.client.envPath(envName, "targetservers", name))
}

// Update updates a target server in the specified environment.
func (s *TargetServerService) Update(ctx context.Context, envName, name string, ts *TargetServer) (*TargetServer, error) {
	if err := ts.Validate(); err != nil {
		return nil, err
	}
	return doUpdate[TargetServer](ctx, s.client, s.client.envPath(envName, "targetservers", name), ts)
}

// Delete deletes a target server from the specified environment.
func (s *TargetServerService) Delete(ctx context.Context, envName, name string) error {
	return doDelete(ctx, s.client, s.client.envPath(envName, "targetservers", name))
}

// List lists all target servers in the specified environment.
// Note: The target server list API does not support pagination and returns only names.
func (s *TargetServerService) List(ctx context.Context, envName string) (*TargetServerListResponse, error) {
	names, err := doListNames(ctx, s.client, s.client.envPath(envName, "targetservers"))
	if err != nil {
		return nil, err
	}
	return &TargetServerListResponse{TargetServerNames: names}, nil
}
