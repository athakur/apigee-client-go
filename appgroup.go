package apigee

import (
	"context"
)

// AppGroup represents an Apigee app group.
type AppGroup struct {
	// Name is the resource name of the app group.
	Name string `json:"name,omitempty"`
	// AppGroupID is the unique identifier of the app group.
	AppGroupID string `json:"appGroupId,omitempty"`
	// DisplayName is the display name of the app group.
	DisplayName string `json:"displayName,omitempty"`
	// ChannelID is the channel identifier.
	ChannelID string `json:"channelId,omitempty"`
	// ChannelURI is the channel URI.
	ChannelURI string `json:"channelUri,omitempty"`
	// Status is the status of the app group (e.g., "active", "inactive").
	Status string `json:"status,omitempty"`
	// Attributes are custom attributes for the app group.
	Attributes []Attribute `json:"attributes,omitempty"`
	// CreatedAt is the creation timestamp in milliseconds (as a string).
	CreatedAt string `json:"createdAt,omitempty"`
	// LastModifiedAt is the last modification timestamp in milliseconds (as a string).
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`
}

// AppGroupListResponse is the response for listing app groups.
type AppGroupListResponse struct {
	// AppGroups is the list of app groups.
	AppGroups []AppGroup `json:"appGroups,omitempty"`
	// NextPageToken is the token for the next page of results.
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// AppGroupService handles operations on app groups.
type AppGroupService struct {
	client *Client
}

// Create creates a new app group.
func (s *AppGroupService) Create(ctx context.Context, appGroup *AppGroup) (*AppGroup, error) {
	return doCreate[AppGroup](ctx, s.client, s.client.orgPath("appgroups"), appGroup)
}

// Get retrieves an app group by name.
func (s *AppGroupService) Get(ctx context.Context, name string) (*AppGroup, error) {
	return doGet[AppGroup](ctx, s.client, s.client.orgPath("appgroups", name))
}

// Update updates an existing app group.
func (s *AppGroupService) Update(ctx context.Context, name string, appGroup *AppGroup) (*AppGroup, error) {
	return doUpdate[AppGroup](ctx, s.client, s.client.orgPath("appgroups", name), appGroup)
}

// Delete deletes an app group.
func (s *AppGroupService) Delete(ctx context.Context, name string) error {
	return doDelete(ctx, s.client, s.client.orgPath("appgroups", name))
}

// List lists all app groups in the organization.
func (s *AppGroupService) List(ctx context.Context, opts *ListOptions) (*AppGroupListResponse, error) {
	return doList[AppGroupListResponse](ctx, s.client, s.client.orgPath("appgroups"), opts)
}
