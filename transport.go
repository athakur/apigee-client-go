package apigee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// buildPath constructs a URL path by joining the base URL with path segments.
func (c *Client) buildPath(segments ...string) string {
	var escaped []string
	for _, s := range segments {
		escaped = append(escaped, url.PathEscape(s))
	}
	return c.BaseURL + "/" + strings.Join(escaped, "/")
}

// newRequest creates a new HTTP request with the given method, URL, and optional body.
func (c *Client) newRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("apigee: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("apigee: failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// doRequest executes an HTTP request and returns the response body.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// Add authentication if token source is available
	if c.tokenSource != nil {
		token, err := c.tokenSource.Token()
		if err != nil {
			return nil, fmt.Errorf("apigee: failed to get token: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("apigee: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("apigee: failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseError(resp, body)
	}

	return body, nil
}

// do executes an HTTP request and unmarshals the response into v.
func (c *Client) do(ctx context.Context, method, urlStr string, reqBody, v interface{}) error {
	req, err := c.newRequest(ctx, method, urlStr, reqBody)
	if err != nil {
		return err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if v != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, v); err != nil {
			return fmt.Errorf("apigee: failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// addQueryParams adds query parameters to a URL string.
func addQueryParams(urlStr string, opts *ListOptions) string {
	if opts == nil {
		return urlStr
	}

	params := url.Values{}
	if opts.PageSize > 0 {
		params.Set("pageSize", fmt.Sprintf("%d", opts.PageSize))
	}
	if opts.PageToken != "" {
		params.Set("pageToken", opts.PageToken)
	}
	if opts.Filter != "" {
		params.Set("filter", opts.Filter)
	}

	if len(params) == 0 {
		return urlStr
	}

	return urlStr + "?" + params.Encode()
}
