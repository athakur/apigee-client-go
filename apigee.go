package apigee

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// DefaultBaseURL is the default base URL for the Apigee API.
	DefaultBaseURL = "https://apigee.googleapis.com/v1"

	// CloudPlatformScope is the OAuth2 scope required for Apigee API access.
	CloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

	// DefaultMaxResponseBytes is the maximum response body size (10 MB).
	DefaultMaxResponseBytes int64 = 10 * 1024 * 1024

	// defaultUserAgent is the User-Agent header sent with all requests.
	defaultUserAgent = "apigee-client-go"

	// DefaultRequestTimeout is the default per-request timeout.
	// It applies when the caller's context has no deadline.
	// Callers can override per-call with context.WithTimeout.
	DefaultRequestTimeout = 60 * time.Second
)

// Client is the Apigee API client.
type Client struct {
	// Organization is the Apigee organization name.
	Organization string

	// BaseURL is the base URL for the Apigee API.
	BaseURL string

	// httpClient is the HTTP client used for requests.
	httpClient *http.Client

	// tokenSource is the OAuth2 token source for authentication.
	tokenSource oauth2.TokenSource

	// maxResponseBytes is the maximum allowed response body size in bytes.
	maxResponseBytes int64

	// userAgent is the User-Agent header sent with all requests.
	userAgent string

	// requestTimeout is the default per-request timeout applied when the
	// caller's context has no deadline. Zero means no timeout.
	requestTimeout time.Duration

	// Services
	AppGroups             *AppGroupService
	AppGroupApps          *AppGroupAppService
	AppGroupAppKeys       *AppGroupAppKeyService
	APIProducts           *APIProductService
	KeyValueMaps          *KeyValueMapService
	KeyValueMapEntries    *KeyValueMapEntryService
	EnvKeyValueMaps       *EnvKeyValueMapService
	EnvKeyValueMapEntries *EnvKeyValueMapEntryService
	TargetServers         *TargetServerService
}

// Attribute represents a name-value attribute.
type Attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// ListOptions specifies pagination and filtering options for list operations.
type ListOptions struct {
	// PageSize is the maximum number of items to return.
	PageSize int
	// PageToken is the token for the next page of results.
	PageToken string
	// Filter is an optional filter expression.
	Filter string
}

// NewClient creates a new Apigee API client.
//
// By default, the client uses Google Application Default Credentials (ADC).
// Use functional options to customize the client configuration.
func NewClient(ctx context.Context, organization string, opts ...ClientOption) (*Client, error) {
	if organization == "" {
		return nil, fmt.Errorf("apigee: organization is required")
	}

	c := &Client{
		Organization:     organization,
		BaseURL:          DefaultBaseURL,
		httpClient:       http.DefaultClient,
		maxResponseBytes: DefaultMaxResponseBytes,
		userAgent:        defaultUserAgent,
		requestTimeout:   DefaultRequestTimeout,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Set up default authentication if no token source or custom HTTP client provided
	if c.tokenSource == nil {
		ts, err := google.DefaultTokenSource(ctx, CloudPlatformScope)
		if err != nil {
			return nil, fmt.Errorf("apigee: failed to create default token source: %w", err)
		}
		c.tokenSource = ts
	}

	// Initialize services
	c.AppGroups = &AppGroupService{client: c}
	c.AppGroupApps = &AppGroupAppService{client: c}
	c.AppGroupAppKeys = &AppGroupAppKeyService{client: c}
	c.APIProducts = &APIProductService{client: c}
	c.KeyValueMaps = &KeyValueMapService{client: c}
	c.KeyValueMapEntries = &KeyValueMapEntryService{client: c}
	c.EnvKeyValueMaps = &EnvKeyValueMapService{client: c}
	c.EnvKeyValueMapEntries = &EnvKeyValueMapEntryService{client: c}
	c.TargetServers = &TargetServerService{client: c}

	return c, nil
}
