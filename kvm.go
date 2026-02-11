package apigee

import (
	"context"
)

// KeyValueMap represents an Apigee Key Value Map (KVM).
type KeyValueMap struct {
	// Name is the KVM identifier.
	Name string `json:"name,omitempty"`
	// Encrypted indicates if the KVM values are encrypted. Always true in Apigee X/hybrid.
	Encrypted bool `json:"encrypted"`
	// MaskedValues indicates if values should be masked in responses.
	MaskedValues bool `json:"maskedValues"`
}

// KeyValueMapListResponse is the response for listing KVMs.
// Note: The Apigee API returns a simple array of KVM names.
type KeyValueMapListResponse struct {
	// KeyValueMapNames is the list of KVM names.
	KeyValueMapNames []string
}

// KeyValueEntry represents an entry in a Key Value Map.
type KeyValueEntry struct {
	// Name is the entry key.
	Name string `json:"name,omitempty"`
	// Value is the entry value.
	Value string `json:"value,omitempty"`
}

// KeyValueEntryListResponse is the response for listing KVM entries.
type KeyValueEntryListResponse struct {
	// KeyValueEntries is the list of KVM entries.
	KeyValueEntries []KeyValueEntry `json:"keyValueEntries,omitempty"`
	// NextPageToken is the token for the next page of results.
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// KeyValueMapService handles operations on organization-level Key Value Maps.
type KeyValueMapService struct {
	client *Client
}

// Create creates a new organization-level KVM.
func (s *KeyValueMapService) Create(ctx context.Context, kvm *KeyValueMap) (*KeyValueMap, error) {
	return doCreate[KeyValueMap](ctx, s.client, s.client.orgPath("keyvaluemaps"), kvm)
}

// Get retrieves an organization-level KVM by name.
func (s *KeyValueMapService) Get(ctx context.Context, name string) (*KeyValueMap, error) {
	return doGet[KeyValueMap](ctx, s.client, s.client.orgPath("keyvaluemaps", name))
}

// Delete deletes an organization-level KVM.
func (s *KeyValueMapService) Delete(ctx context.Context, name string) error {
	return doDelete(ctx, s.client, s.client.orgPath("keyvaluemaps", name))
}

// List lists all organization-level KVMs.
// Note: The KVM list API does not support pagination and returns only KVM names.
func (s *KeyValueMapService) List(ctx context.Context) (*KeyValueMapListResponse, error) {
	names, err := doListNames(ctx, s.client, s.client.orgPath("keyvaluemaps"))
	if err != nil {
		return nil, err
	}
	return &KeyValueMapListResponse{KeyValueMapNames: names}, nil
}

// KeyValueMapEntryService handles operations on organization-level KVM entries.
type KeyValueMapEntryService struct {
	client *Client
}

// Create creates a new entry in an organization-level KVM.
func (s *KeyValueMapEntryService) Create(ctx context.Context, kvmName string, entry *KeyValueEntry) (*KeyValueEntry, error) {
	return doCreate[KeyValueEntry](ctx, s.client, s.client.orgPath("keyvaluemaps", kvmName, "entries"), entry)
}

// Get retrieves an entry from an organization-level KVM.
func (s *KeyValueMapEntryService) Get(ctx context.Context, kvmName, entryName string) (*KeyValueEntry, error) {
	return doGet[KeyValueEntry](ctx, s.client, s.client.orgPath("keyvaluemaps", kvmName, "entries", entryName))
}

// Update updates an entry in an organization-level KVM.
func (s *KeyValueMapEntryService) Update(ctx context.Context, kvmName, entryName string, entry *KeyValueEntry) (*KeyValueEntry, error) {
	return doUpdate[KeyValueEntry](ctx, s.client, s.client.orgPath("keyvaluemaps", kvmName, "entries", entryName), entry)
}

// Delete deletes an entry from an organization-level KVM.
func (s *KeyValueMapEntryService) Delete(ctx context.Context, kvmName, entryName string) error {
	return doDelete(ctx, s.client, s.client.orgPath("keyvaluemaps", kvmName, "entries", entryName))
}

// List lists all entries in an organization-level KVM.
// Note: The KVM entries list API does not support pagination.
func (s *KeyValueMapEntryService) List(ctx context.Context, kvmName string) (*KeyValueEntryListResponse, error) {
	return doGet[KeyValueEntryListResponse](ctx, s.client, s.client.orgPath("keyvaluemaps", kvmName, "entries"))
}

// EnvKeyValueMapService handles operations on environment-level Key Value Maps.
type EnvKeyValueMapService struct {
	client *Client
}

// Create creates a new environment-level KVM.
func (s *EnvKeyValueMapService) Create(ctx context.Context, envName string, kvm *KeyValueMap) (*KeyValueMap, error) {
	return doCreate[KeyValueMap](ctx, s.client, s.client.envPath(envName, "keyvaluemaps"), kvm)
}

// Get retrieves an environment-level KVM by name.
func (s *EnvKeyValueMapService) Get(ctx context.Context, envName, name string) (*KeyValueMap, error) {
	return doGet[KeyValueMap](ctx, s.client, s.client.envPath(envName, "keyvaluemaps", name))
}

// Delete deletes an environment-level KVM.
func (s *EnvKeyValueMapService) Delete(ctx context.Context, envName, name string) error {
	return doDelete(ctx, s.client, s.client.envPath(envName, "keyvaluemaps", name))
}

// List lists all environment-level KVMs.
// Note: The KVM list API does not support pagination and returns only KVM names.
func (s *EnvKeyValueMapService) List(ctx context.Context, envName string) (*KeyValueMapListResponse, error) {
	names, err := doListNames(ctx, s.client, s.client.envPath(envName, "keyvaluemaps"))
	if err != nil {
		return nil, err
	}
	return &KeyValueMapListResponse{KeyValueMapNames: names}, nil
}

// EnvKeyValueMapEntryService handles operations on environment-level KVM entries.
type EnvKeyValueMapEntryService struct {
	client *Client
}

// Create creates a new entry in an environment-level KVM.
func (s *EnvKeyValueMapEntryService) Create(ctx context.Context, envName, kvmName string, entry *KeyValueEntry) (*KeyValueEntry, error) {
	return doCreate[KeyValueEntry](ctx, s.client, s.client.envPath(envName, "keyvaluemaps", kvmName, "entries"), entry)
}

// Get retrieves an entry from an environment-level KVM.
func (s *EnvKeyValueMapEntryService) Get(ctx context.Context, envName, kvmName, entryName string) (*KeyValueEntry, error) {
	return doGet[KeyValueEntry](ctx, s.client, s.client.envPath(envName, "keyvaluemaps", kvmName, "entries", entryName))
}

// Update updates an entry in an environment-level KVM.
func (s *EnvKeyValueMapEntryService) Update(ctx context.Context, envName, kvmName, entryName string, entry *KeyValueEntry) (*KeyValueEntry, error) {
	return doUpdate[KeyValueEntry](ctx, s.client, s.client.envPath(envName, "keyvaluemaps", kvmName, "entries", entryName), entry)
}

// Delete deletes an entry from an environment-level KVM.
func (s *EnvKeyValueMapEntryService) Delete(ctx context.Context, envName, kvmName, entryName string) error {
	return doDelete(ctx, s.client, s.client.envPath(envName, "keyvaluemaps", kvmName, "entries", entryName))
}

// List lists all entries in an environment-level KVM.
// Note: The KVM entries list API does not support pagination.
func (s *EnvKeyValueMapEntryService) List(ctx context.Context, envName, kvmName string) (*KeyValueEntryListResponse, error) {
	return doGet[KeyValueEntryListResponse](ctx, s.client, s.client.envPath(envName, "keyvaluemaps", kvmName, "entries"))
}
