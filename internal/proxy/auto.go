// Package proxy provides automatic proxy retry functionality.
package proxy

import (
	"context"
	"net/http"

	"github.com/tritueviet/search-agents/internal/httpclient"
)

// AutoProxyClient wraps HTTP client with automatic proxy retry.
type AutoProxyClient struct {
	client     *httpclient.Client
	torProxy   string
	httpProxy  string
}

// NewAutoProxyClient creates a new auto-proxy client.
func NewAutoProxyClient(client *httpclient.Client) *AutoProxyClient {
	return &AutoProxyClient{
		client:   client,
		torProxy: "socks5h://127.0.0.1:9150",
		httpProxy: "",
	}
}

// SetTorProxy sets custom Tor proxy URL.
func (a *AutoProxyClient) SetTorProxy(url string) {
	a.torProxy = url
}

// SetHTTPProxy sets custom HTTP proxy URL.
func (a *AutoProxyClient) SetHTTPProxy(url string) {
	a.httpProxy = url
}

// DoWithProxyRetry executes request with automatic proxy retry.
func (a *AutoProxyClient) DoWithProxyRetry(ctx context.Context, method, url string, body interface{}) (*httpclient.Response, error) {
	// First attempt with normal client
	resp, err := a.client.Do(ctx, method, url, nil)
	if err == nil && resp.StatusCode != 403 {
		return resp, nil
	}

	// If 403 or error, try Tor proxy
	if a.torProxy != "" {
		torClient, torErr := httpclient.NewClient(httpclient.Options{
			Proxy:   a.torProxy,
			Timeout: 15,
			Verify:  true,
		})
		if torErr == nil {
			torClient.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			torResp, torRespErr := torClient.Do(ctx, method, url, nil)
			if torRespErr == nil && torResp.StatusCode != 403 {
				return torResp, nil
			}
		}
	}

	// If Tor fails, try HTTP proxy if configured
	if a.httpProxy != "" {
		httpProxyClient, proxyErr := httpclient.NewClient(httpclient.Options{
			Proxy:   a.httpProxy,
			Timeout: 15,
			Verify:  true,
		})
		if proxyErr == nil {
			httpProxyClient.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			proxyResp, proxyRespErr := httpProxyClient.Do(ctx, method, url, nil)
			if proxyRespErr == nil && proxyResp.StatusCode != 403 {
				return proxyResp, nil
			}
		}
	}

	// Return original response/error
	if resp != nil {
		return resp, err
	}

	return nil, err
}

// Get creates a GET request with proxy retry.
func (a *AutoProxyClient) Get(ctx context.Context, url string) (*httpclient.Response, error) {
	return a.DoWithProxyRetry(ctx, http.MethodGet, url, nil)
}

// Post creates a POST request with proxy retry.
func (a *AutoProxyClient) Post(ctx context.Context, url string, contentType string, body interface{}) (*httpclient.Response, error) {
	return a.DoWithProxyRetry(ctx, http.MethodPost, url, body)
}
