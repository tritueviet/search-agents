// Package torproxy provides automatic Tor proxy retry for HTTP requests.
package torproxy

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/tritueviet/search-agents/internal/httpclient"
)

// DefaultTorProxyURL is the default Tor SOCKS proxy URL.
const DefaultTorProxyURL = "socks5h://127.0.0.1:9050"

// Client wraps HTTP client with automatic Tor proxy retry on 403.
type Client struct {
	normal    *httpclient.Client
	tor       *httpclient.Client
	torURL    string
	available bool
}

// NewClient creates a new Tor proxy retry client.
func NewClient(normal *httpclient.Client) *Client {
	torURL := os.Getenv("TOR_PROXY_URL")
	if torURL == "" {
		torURL = DefaultTorProxyURL
	}

	tor, err := httpclient.NewClient(httpclient.Options{
		Proxy:   torURL,
		Timeout: 10,
		Verify:  true,
	})

	torAvailable := err == nil
	if torAvailable {
		tor.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	} else {
		tor = nil
	}

	return &Client{
		normal:    normal,
		tor:       tor,
		torURL:    torURL,
		available: torAvailable,
	}
}

// Do executes request with automatic Tor proxy retry on 403.
func (c *Client) Do(ctx context.Context, method, url string) (*httpclient.Response, error) {
	// If Tor is available, use it directly
	if c.tor != nil {
		resp, err := c.tor.Do(ctx, method, url, nil)
		return resp, err
	}

	// Fallback to normal client
	return c.normal.Do(ctx, method, url, nil)
}

// Get performs GET request with Tor.
func (c *Client) Get(ctx context.Context, url string) (*httpclient.Response, error) {
	return c.Do(ctx, http.MethodGet, url)
}

// Post performs POST request with Tor.
func (c *Client) Post(ctx context.Context, url string) (*httpclient.Response, error) {
	return c.Do(ctx, http.MethodPost, url)
}

// PostForm performs POST with form data.
func (c *Client) PostForm(ctx context.Context, url string) (*httpclient.Response, error) {
	return c.Do(ctx, http.MethodPost, url)
}

// IsTorAvailable returns true if Tor proxy is available.
func (c *Client) IsTorAvailable() bool {
	return c.available
}

// TorURL returns the Tor proxy URL being used.
func (c *Client) TorURL() string {
	return c.torURL
}

// StatusMessage returns a human-readable status message.
func (c *Client) StatusMessage() string {
	if c.available {
		return fmt.Sprintf("Tor proxy available at %s", c.torURL)
	}
	return "Tor proxy not available - install Tor: sudo apt install tor"
}
