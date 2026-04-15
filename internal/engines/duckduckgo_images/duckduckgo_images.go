// Package duckduckgo_images implements DuckDuckGo image search.
package duckduckgo_images

import (
	"context"
	"fmt"
	"strings"

	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
)

// DuckduckgoImages implements DuckDuckGo image search engine.
type DuckduckgoImages struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

// New creates a new DuckDuckGo images engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &DuckduckgoImages{
		client:   client,
		name:     "duckduckgo_images",
		category: engine.CategoryImages,
		provider: "duckduckgo",
		priority: 0.0,
	}
}

func (d *DuckduckgoImages) Name() string         { return d.name }
func (d *DuckduckgoImages) Category() engine.Category { return d.category }
func (d *DuckduckgoImages) Provider() string     { return d.provider }
func (d *DuckduckgoImages) Priority() float64    { return d.priority }

// Search performs DuckDuckGo image search.
// NOTE: DuckDuckGo JSON API (/i.js) has been deprecated and returns 403.
func (d *DuckduckgoImages) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	return nil, fmt.Errorf("DuckDuckGo Images API has been deprecated (403 Forbidden). " +
		"Alternatives:\n" +
		"  1. Use Bing Images API (requires API key)\n" +
		"  2. Use Google Custom Search API (requires API key)\n" +
		"  3. Use SerpAPI or similar third-party services")
}

func (d *DuckduckgoImages) buildPayload(query string, opts engine.SearchOptions) string {
	return strings.ReplaceAll(query, " ", "+")
}
