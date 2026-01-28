package apigee

import "context"

// APIProductClient defines the interface for API Product operations.
type APIProductClient interface {
	Create(ctx context.Context, product *APIProduct) (*APIProduct, error)
	Get(ctx context.Context, name string) (*APIProduct, error)
	Update(ctx context.Context, name string, product *APIProduct) (*APIProduct, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *ListOptions) (*APIProductListResponse, error)
}

// AppGroupClient defines the interface for App Group operations.
type AppGroupClient interface {
	Create(ctx context.Context, appGroup *AppGroup) (*AppGroup, error)
	Get(ctx context.Context, name string) (*AppGroup, error)
	Update(ctx context.Context, name string, appGroup *AppGroup) (*AppGroup, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *ListOptions) (*AppGroupListResponse, error)
}

// AppGroupAppClient defines the interface for App Group App operations.
type AppGroupAppClient interface {
	Create(ctx context.Context, appGroupName string, app *AppGroupApp) (*AppGroupApp, error)
	Get(ctx context.Context, appGroupName, appName string) (*AppGroupApp, error)
	Update(ctx context.Context, appGroupName, appName string, app *AppGroupApp) (*AppGroupApp, error)
	Delete(ctx context.Context, appGroupName, appName string) error
	List(ctx context.Context, appGroupName string, opts *ListOptions) (*AppGroupAppListResponse, error)
}

// AppGroupAppKeyClient defines the interface for App Group App Key operations.
type AppGroupAppKeyClient interface {
	Create(ctx context.Context, appGroupName, appName string, key *AppGroupAppKeyCreateRequest) (*AppGroupAppKey, error)
	Get(ctx context.Context, appGroupName, appName, consumerKey string) (*AppGroupAppKey, error)
	Update(ctx context.Context, appGroupName, appName, consumerKey string, key *AppGroupAppKey) (*AppGroupAppKey, error)
	Delete(ctx context.Context, appGroupName, appName, consumerKey string) error
	UpdateAPIProductStatus(ctx context.Context, appGroupName, appName, consumerKey, apiProduct, action string) error
	Generate(ctx context.Context, appGroupName, appName string, req *AppGroupAppKeyGenerateRequest) (*AppGroupApp, error)
}

// KeyValueMapClient defines the interface for organization-level KVM operations.
type KeyValueMapClient interface {
	Create(ctx context.Context, kvm *KeyValueMap) (*KeyValueMap, error)
	Get(ctx context.Context, name string) (*KeyValueMap, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) (*KeyValueMapListResponse, error)
}

// KeyValueMapEntryClient defines the interface for organization-level KVM entry operations.
type KeyValueMapEntryClient interface {
	Create(ctx context.Context, kvmName string, entry *KeyValueEntry) (*KeyValueEntry, error)
	Get(ctx context.Context, kvmName, entryName string) (*KeyValueEntry, error)
	Update(ctx context.Context, kvmName, entryName string, entry *KeyValueEntry) (*KeyValueEntry, error)
	Delete(ctx context.Context, kvmName, entryName string) error
	List(ctx context.Context, kvmName string) (*KeyValueEntryListResponse, error)
}

// EnvKeyValueMapClient defines the interface for environment-level KVM operations.
type EnvKeyValueMapClient interface {
	Create(ctx context.Context, envName string, kvm *KeyValueMap) (*KeyValueMap, error)
	Get(ctx context.Context, envName, name string) (*KeyValueMap, error)
	Delete(ctx context.Context, envName, name string) error
	List(ctx context.Context, envName string) (*KeyValueMapListResponse, error)
}

// EnvKeyValueMapEntryClient defines the interface for environment-level KVM entry operations.
type EnvKeyValueMapEntryClient interface {
	Create(ctx context.Context, envName, kvmName string, entry *KeyValueEntry) (*KeyValueEntry, error)
	Get(ctx context.Context, envName, kvmName, entryName string) (*KeyValueEntry, error)
	Update(ctx context.Context, envName, kvmName, entryName string, entry *KeyValueEntry) (*KeyValueEntry, error)
	Delete(ctx context.Context, envName, kvmName, entryName string) error
	List(ctx context.Context, envName, kvmName string) (*KeyValueEntryListResponse, error)
}

// TargetServerClient defines the interface for Target Server operations.
type TargetServerClient interface {
	Create(ctx context.Context, envName string, ts *TargetServer) (*TargetServer, error)
	Get(ctx context.Context, envName, name string) (*TargetServer, error)
	Update(ctx context.Context, envName, name string, ts *TargetServer) (*TargetServer, error)
	Delete(ctx context.Context, envName, name string) error
	List(ctx context.Context, envName string) (*TargetServerListResponse, error)
}

// Compile-time interface implementation checks.
var (
	_ APIProductClient         = (*APIProductService)(nil)
	_ AppGroupClient           = (*AppGroupService)(nil)
	_ AppGroupAppClient        = (*AppGroupAppService)(nil)
	_ AppGroupAppKeyClient     = (*AppGroupAppKeyService)(nil)
	_ KeyValueMapClient        = (*KeyValueMapService)(nil)
	_ KeyValueMapEntryClient   = (*KeyValueMapEntryService)(nil)
	_ EnvKeyValueMapClient     = (*EnvKeyValueMapService)(nil)
	_ EnvKeyValueMapEntryClient = (*EnvKeyValueMapEntryService)(nil)
	_ TargetServerClient       = (*TargetServerService)(nil)
)
