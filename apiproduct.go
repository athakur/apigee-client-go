package apigee

import (
	"context"
	"net/http"
)

// APIProduct represents an Apigee API product.
type APIProduct struct {
	// Name is the name of the API product.
	Name string `json:"name,omitempty"`
	// DisplayName is the display name of the API product.
	DisplayName string `json:"displayName,omitempty"`
	// Description is the description of the API product.
	Description string `json:"description,omitempty"`
	// ApprovalType is the approval type ("auto" or "manual").
	ApprovalType string `json:"approvalType,omitempty"`
	// Proxies is the list of API proxies associated with the product.
	Proxies []string `json:"proxies,omitempty"`
	// Environments is the list of environments where the product is available.
	Environments []string `json:"environments,omitempty"`
	// APIResources is the list of API resource paths included in the product.
	APIResources []string `json:"apiResources,omitempty"`
	// Scopes is the list of OAuth scopes for the product.
	Scopes []string `json:"scopes,omitempty"`
	// Quota is the quota limit for the product.
	Quota string `json:"quota,omitempty"`
	// QuotaInterval is the quota interval.
	QuotaInterval string `json:"quotaInterval,omitempty"`
	// QuotaTimeUnit is the quota time unit (e.g., "minute", "hour", "day", "month").
	QuotaTimeUnit string `json:"quotaTimeUnit,omitempty"`
	// Attributes are custom attributes for the product.
	Attributes []Attribute `json:"attributes,omitempty"`
	// CreatedAt is the creation timestamp in milliseconds (as a string).
	CreatedAt string `json:"createdAt,omitempty"`
	// LastModifiedAt is the last modification timestamp in milliseconds (as a string).
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`
}

// APIProductListResponse is the response for listing API products.
type APIProductListResponse struct {
	// APIProducts is the list of API products.
	APIProducts []APIProduct `json:"apiProduct,omitempty"`
	// NextPageToken is the token for the next page of results.
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// APIProductService handles operations on API products.
type APIProductService struct {
	client *Client
}

// Create creates a new API product.
func (s *APIProductService) Create(ctx context.Context, product *APIProduct) (*APIProduct, error) {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "apiproducts")

	result := &APIProduct{}
	if err := s.client.do(ctx, http.MethodPost, endpoint, product, result); err != nil {
		return nil, err
	}

	return result, nil
}

// Get retrieves an API product by name.
func (s *APIProductService) Get(ctx context.Context, name string) (*APIProduct, error) {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "apiproducts", name)

	result := &APIProduct{}
	if err := s.client.do(ctx, http.MethodGet, endpoint, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates an existing API product.
func (s *APIProductService) Update(ctx context.Context, name string, product *APIProduct) (*APIProduct, error) {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "apiproducts", name)

	result := &APIProduct{}
	if err := s.client.do(ctx, http.MethodPut, endpoint, product, result); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete deletes an API product.
func (s *APIProductService) Delete(ctx context.Context, name string) error {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "apiproducts", name)

	return s.client.do(ctx, http.MethodDelete, endpoint, nil, nil)
}

// List lists all API products in the organization.
func (s *APIProductService) List(ctx context.Context, opts *ListOptions) (*APIProductListResponse, error) {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "apiproducts")
	endpoint = addQueryParams(endpoint, opts)

	result := &APIProductListResponse{}
	if err := s.client.do(ctx, http.MethodGet, endpoint, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}
