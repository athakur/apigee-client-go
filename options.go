package apigee

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client for the Apigee client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL for the Apigee API.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithTokenSource sets a custom OAuth2 token source for authentication.
func WithTokenSource(ts oauth2.TokenSource) ClientOption {
	return func(c *Client) {
		c.tokenSource = ts
	}
}

// WithMaxResponseBytes sets the maximum allowed response body size in bytes.
// The default is DefaultMaxResponseBytes (10 MB). If the response exceeds this
// limit, the request will fail with an error.
func WithMaxResponseBytes(n int64) ClientOption {
	return func(c *Client) {
		c.maxResponseBytes = n
	}
}

// WithUserAgent sets the User-Agent header sent with all requests.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// WithRequestTimeout sets the default per-request timeout. It applies when the
// caller's context has no deadline. Callers can still override individual calls
// by passing a context with a shorter deadline via context.WithTimeout.
// The default is DefaultRequestTimeout (60s). Set to 0 to disable.
func WithRequestTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.requestTimeout = d
	}
}
