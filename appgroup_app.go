package apigee

import (
	"context"
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
	return doCreate[AppGroupApp](ctx, s.client, s.client.orgPath("appgroups", appGroupName, "apps"), app)
}

// Get retrieves an app by name from an app group.
func (s *AppGroupAppService) Get(ctx context.Context, appGroupName, appName string) (*AppGroupApp, error) {
	return doGet[AppGroupApp](ctx, s.client, s.client.orgPath("appgroups", appGroupName, "apps", appName))
}

// Update updates an existing app in an app group.
func (s *AppGroupAppService) Update(ctx context.Context, appGroupName, appName string, app *AppGroupApp) (*AppGroupApp, error) {
	return doUpdate[AppGroupApp](ctx, s.client, s.client.orgPath("appgroups", appGroupName, "apps", appName), app)
}

// Delete deletes an app from an app group.
func (s *AppGroupAppService) Delete(ctx context.Context, appGroupName, appName string) error {
	return doDelete(ctx, s.client, s.client.orgPath("appgroups", appGroupName, "apps", appName))
}

// List lists all apps in an app group.
func (s *AppGroupAppService) List(ctx context.Context, appGroupName string, opts *ListOptions) (*AppGroupAppListResponse, error) {
	return doList[AppGroupAppListResponse](ctx, s.client, s.client.orgPath("appgroups", appGroupName, "apps"), opts)
}
