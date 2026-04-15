// Package httpclient provides an HTTP client with proxy and timeout support.
package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Response wraps an HTTP response with helper methods.
type Response struct {
	StatusCode int
	Content    []byte
	Text       string
	Headers    http.Header
}

// Client is an HTTP client with proxy and timeout support.
type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// Options contains configuration for the HTTP client.
type Options struct {
	Proxy   string
	Timeout time.Duration
	Verify  bool
}

// NewClient creates a new HTTP client.
func NewClient(opts Options) (*Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !opts.Verify, //nolint:gosec
		},
	}

	// Set up proxy if provided
	if opts.Proxy != "" {
		proxyURL, err := url.Parse(opts.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	return &Client{
		client:  client,
		headers: make(map[string]string),
	}, nil
}

// SetHeader sets a header for all requests.
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// SetHeaders sets multiple headers for all requests.
func (c *Client) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		c.headers[k] = v
	}
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, rawURL string) (*Response, error) {
	return c.Do(ctx, http.MethodGet, rawURL, nil)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, rawURL string, contentType string, body io.Reader) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	return c.doRequest(req)
}

// PostForm performs a POST request with form data.
func (c *Client) PostForm(ctx context.Context, rawURL string, data url.Values) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}

// Do performs an HTTP request.
func (c *Client) Do(ctx context.Context, method, rawURL string, body io.Reader) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return c.doRequest(req)
}

// doRequest executes a request with the configured headers.
func (c *Client) doRequest(req *http.Request) (*Response, error) {
	// Set custom headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Set common headers to mimic a real browser
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	}
	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		// Check for timeout
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			return nil, fmt.Errorf("request timed out: %w", err)
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Content:    content,
		Text:       string(content),
		Headers:    resp.Header,
	}, nil
}
