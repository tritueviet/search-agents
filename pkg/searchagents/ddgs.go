// Package searchagents provides the public API for Search Agents.
package searchagents

import (
	"context"

	"github.com/tritueviet/search-agents/internal/core"
	"github.com/tritueviet/search-agents/internal/engine"
)

// SearchAgents is the main entry point for metasearch operations.
type SearchAgents struct {
	core *core.SearchAgents
}

// Options contains configuration for SearchAgents.
type Options struct {
	// Proxy is the proxy URL (e.g., "socks5h://127.0.0.1:9150").
	Proxy string
	// Timeout is the timeout in seconds for HTTP requests.
	Timeout int
	// Verify enables SSL verification (default: true).
	Verify bool
}

// New creates a new SearchAgents instance.
func New(opts Options) (*SearchAgents, error) {
	c, err := core.New(core.Options{
		Proxy:   opts.Proxy,
		Timeout: opts.Timeout,
		Verify:  opts.Verify,
	})
	if err != nil {
		return nil, err
	}

	return &SearchAgents{core: c}, nil
}

// Text performs a text search.
func (s *SearchAgents) Text(ctx context.Context, query string, opts ...engine.SearchOptions) ([]map[string]string, error) {
	var opt engine.SearchOptions
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = engine.DefaultSearchOptions()
	}

	return s.core.Text(ctx, query, opt)
}

// Images performs an image search.
func (s *SearchAgents) Images(ctx context.Context, query string, opts ...engine.SearchOptions) ([]map[string]string, error) {
	var opt engine.SearchOptions
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = engine.DefaultSearchOptions()
	}

	return s.core.Images(ctx, query, opt)
}

// Videos performs a video search.
func (s *SearchAgents) Videos(ctx context.Context, query string, opts ...engine.SearchOptions) ([]map[string]string, error) {
	var opt engine.SearchOptions
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = engine.DefaultSearchOptions()
	}

	return s.core.Videos(ctx, query, opt)
}

// News performs a news search.
func (s *SearchAgents) News(ctx context.Context, query string, opts ...engine.SearchOptions) ([]map[string]string, error) {
	var opt engine.SearchOptions
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = engine.DefaultSearchOptions()
	}

	return s.core.News(ctx, query, opt)
}

// Books performs a books search.
func (s *SearchAgents) Books(ctx context.Context, query string, opts ...engine.SearchOptions) ([]map[string]string, error) {
	var opt engine.SearchOptions
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = engine.DefaultSearchOptions()
	}

	return s.core.Books(ctx, query, opt)
}
