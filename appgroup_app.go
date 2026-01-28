package apigee

import (
	"context"
	"net/http"
)

// AppGroupApp represents an app within an app group.
type AppGroupApp struct {
	// Name is the name of the app.
	Name string `json:"name,omitempty"`
	// AppID is the unique identifier of the app.
	AppID string `json:"appId,omitempty"`
	// AppGroup is the name of the parent app group.
	AppGroup string `json:"appGroup,omitempty"`
	// Status is the status of the app (e.g., "approved", "revoked").
	Status string `json:"status,omitempty"`
	// CallbackURL is the callback URL for OAuth.
	CallbackURL string `json:"callbackUrl,omitempty"`
	// APIProducts is the list of API products associated with the app.
	APIProducts []string `json:"apiProducts,omitempty"`
	// Credentials is the list of credentials (keys) for the app.
	Credentials []AppGroupAppKey `json:"credentials,omitempty"`
	// Attributes are custom attributes for the app.
	Attributes []Attribute `json:"attributes,omitempty"`
	// Scopes is the list of OAuth scopes for the app.
	Scopes []string `json:"scopes,omitempty"`
	// CreatedAt is the creation timestamp in milliseconds (as a string).
	CreatedAt string `json:"createdAt,omitempty"`
	// LastModifiedAt is the last modification timestamp in milliseconds (as a string).
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`
}

// AppGroupAppListResponse is the response for listing apps in an app group.
type AppGroupAppListResponse struct {
	// AppGroupApps is the list of apps.
	AppGroupApps []AppGroupApp `json:"appGroupApps,omitempty"`
	// NextPageToken is the token for the next page of results.
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// AppGroupAppService handles operations on apps within app groups.
type AppGroupAppService struct {
	client *Client
}

// Create creates a new app in an app group.
func (s *AppGroupAppService) Create(ctx context.Context, appGroupName string, app *AppGroupApp) (*AppGroupApp, error) {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps")

	result := &AppGroupApp{}
	if err := s.client.do(ctx, http.MethodPost, endpoint, app, result); err != nil {
		return nil, err
	}

	return result, nil
}

// Get retrieves an app by name from an app group.
func (s *AppGroupAppService) Get(ctx context.Context, appGroupName, appName string) (*AppGroupApp, error) {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps", appName)

	result := &AppGroupApp{}
	if err := s.client.do(ctx, http.MethodGet, endpoint, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates an existing app in an app group.
func (s *AppGroupAppService) Update(ctx context.Context, appGroupName, appName string, app *AppGroupApp) (*AppGroupApp, error) {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps", appName)

	result := &AppGroupApp{}
	if err := s.client.do(ctx, http.MethodPut, endpoint, app, result); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete deletes an app from an app group.
func (s *AppGroupAppService) Delete(ctx context.Context, appGroupName, appName string) error {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps", appName)

	return s.client.do(ctx, http.MethodDelete, endpoint, nil, nil)
}

// List lists all apps in an app group.
func (s *AppGroupAppService) List(ctx context.Context, appGroupName string, opts *ListOptions) (*AppGroupAppListResponse, error) {
	endpoint := s.client.buildPath("organizations", s.client.Organization, "appgroups", appGroupName, "apps")
	endpoint = addQueryParams(endpoint, opts)

	result := &AppGroupAppListResponse{}
	if err := s.client.do(ctx, http.MethodGet, endpoint, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}
