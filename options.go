package apigee

import (
	"net/http"

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
