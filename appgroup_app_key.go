package apigee

import (
	"context"
	"net/http"
)

// AppGroupAppKey represents a credential (key) for an app.
type AppGroupAppKey struct {
	// ConsumerKey is the consumer key (API key).
	ConsumerKey string `json:"consumerKey,omitempty"`
	// ConsumerSecret is the consumer secret.
	ConsumerSecret string `json:"consumerSecret,omitempty"`
	// Status is the status of the key (e.g., "approved", "revoked").
	Status string `json:"status,omitempty"`
	// APIProducts is the list of API products associated with the key.
	APIProducts []APIProductRef `json:"apiProducts,omitempty"`
	// Attributes are custom attributes for the key.
	Attributes []Attribute `json:"attributes,omitempty"`
	// IssuedAt is the timestamp when the key was issued, in milliseconds (as a string).
	IssuedAt string `json:"issuedAt,omitempty"`
	// ExpiresAt is the timestamp when the key expires, in milliseconds (as a string). "-1" means never expires.
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// APIProductRef represents a reference to an API product with its approval status.
type APIProductRef struct {
	// APIProduct is the name of the API product.
	APIProduct string `json:"apiproduct,omitempty"`
	// Status is the approval status for this product.
	Status string `json:"status,omitempty"`
}

// AppGroupAppKeyCreateRequest is the request body for creating a new key with a custom key/secret.
type AppGroupAppKeyCreateRequest struct {
	// ConsumerKey is the custom consumer key.
	ConsumerKey string `json:"consumerKey,omitempty"`
	// ConsumerSecret is the custom consumer secret.
	ConsumerSecret string `json:"consumerSecret,omitempty"`
	// ExpiresInSeconds is the expiration time in seconds. -1 means never expires.
	ExpiresInSeconds int64 `json:"expiresInSeconds,omitempty"`
}

// AppGroupAppKeyGenerateRequest is the request body for generating a new key where Apigee creates the key/secret.
type AppGroupAppKeyGenerateRequest struct {
	// APIProducts is the list of API product names to associate with the generated key.
	APIProducts []string `json:"apiProducts"`
	// KeyExpiresIn is the expiration time in milliseconds. -1 means never expires.
	KeyExpiresIn int64 `json:"keyExpiresIn,omitempty"`
}

// appGroupAppUpdateRequest is the internal request body for PUT operations on an app.
// It contains only mutable fields to avoid overwriting immutable/output-only fields.
type appGroupAppUpdateRequest struct {
	Attributes  []Attribute `json:"attributes,omitempty"`
	CallbackURL string      `json:"callbackUrl,omitempty"`
	Scopes      []string    `json:"scopes,omitempty"`
	Status      string      `json:"status,omitempty"`
	APIProducts []string    `json:"apiProducts"`
	KeyExpiresIn int64      `json:"keyExpiresIn,omitempty"`
}

// AppGroupAppKeyService handles operations on app credentials (keys).
type AppGroupAppKeyService struct {
	client *Client
}

// Create creates a new key for an app.
func (s *AppGroupAppKeyService) Create(ctx context.Context, appGroupName, appName string, key *AppGroupAppKeyCreateRequest) (*AppGroupAppKey, error) {
	url := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps", appName, "keys")

	result := &AppGroupAppKey{}
	if err := s.client.do(ctx, http.MethodPost, url, key, result); err != nil {
		return nil, err
	}

	return result, nil
}

// Get retrieves a key by consumer key.
func (s *AppGroupAppKeyService) Get(ctx context.Context, appGroupName, appName, consumerKey string) (*AppGroupAppKey, error) {
	url := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps", appName, "keys", consumerKey)

	result := &AppGroupAppKey{}
	if err := s.client.do(ctx, http.MethodGet, url, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates an existing key. This can be used to approve/revoke a key or update its API products.
func (s *AppGroupAppKeyService) Update(ctx context.Context, appGroupName, appName, consumerKey string, key *AppGroupAppKey) (*AppGroupAppKey, error) {
	url := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps", appName, "keys", consumerKey)

	result := &AppGroupAppKey{}
	if err := s.client.do(ctx, http.MethodPut, url, key, result); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete deletes a key.
func (s *AppGroupAppKeyService) Delete(ctx context.Context, appGroupName, appName, consumerKey string) error {
	url := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps", appName, "keys", consumerKey)

	return s.client.do(ctx, http.MethodDelete, url, nil, nil)
}

// UpdateAPIProductStatus updates the approval status of an API product for a specific key.
func (s *AppGroupAppKeyService) UpdateAPIProductStatus(ctx context.Context, appGroupName, appName, consumerKey, apiProduct, action string) error {
	url := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps", appName, "keys", consumerKey, "apiproducts", apiProduct)
	url = url + "?action=" + action

	return s.client.do(ctx, http.MethodPost, url, nil, nil)
}

// Generate creates a new key for an app where Apigee generates the key/secret.
// This fetches the existing app, preserves its mutable fields, and performs a PUT
// operation with the specified API products and expiration to generate a new key.
// Returns the updated app containing the newly generated credential in its Credentials field.
func (s *AppGroupAppKeyService) Generate(ctx context.Context, appGroupName, appName string, req *AppGroupAppKeyGenerateRequest) (*AppGroupApp, error) {
	// First, fetch the existing app to preserve mutable fields
	existingApp, err := s.client.AppGroupApps.Get(ctx, appGroupName, appName)
	if err != nil {
		return nil, err
	}

	// Build request with mutable fields from existing app plus new key settings
	updateReq := &appGroupAppUpdateRequest{
		Attributes:   existingApp.Attributes,
		CallbackURL:  existingApp.CallbackURL,
		Scopes:       existingApp.Scopes,
		Status:       existingApp.Status,
		APIProducts:  req.APIProducts,
		KeyExpiresIn: req.KeyExpiresIn,
	}

	url := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps", appName)

	result := &AppGroupApp{}
	if err := s.client.do(ctx, http.MethodPut, url, updateReq, result); err != nil {
		return nil, err
	}

	return result, nil
}
