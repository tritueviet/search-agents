// Package wikipedia implements Wikipedia search engine.
package wikipedia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/utils"
)

// Wikipedia implements Wikipedia search engine.
type Wikipedia struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	wikipediaAPIURL = "https://en.wikipedia.org/w/api.php"
)

// New creates a new Wikipedia engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &Wikipedia{
		client:   client,
		name:     "wikipedia",
		category: engine.CategoryText,
		provider: "wikipedia",
		priority: 0.85,
	}
}

func (w *Wikipedia) Name() string         { return w.name }
func (w *Wikipedia) Category() engine.Category { return w.category }
func (w *Wikipedia) Provider() string     { return w.provider }
func (w *Wikipedia) Priority() float64    { return w.priority }

// Search performs Wikipedia search.
func (w *Wikipedia) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("list", "search")
	params.Set("srsearch", query)
	params.Set("format", "json")
	params.Set("srlimit", "10")

	apiURL := fmt.Sprintf("%s?%s", wikipediaAPIURL, params.Encode())

	resp, err := w.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Query struct {
			Search []struct {
				Title   string `json:"title"`
				Snippet string `json:"snippet"`
			} `json:"search"`
		} `json:"query"`
	}

	if err := json.Unmarshal(resp.Content, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	results := make([]map[string]string, 0, len(result.Query.Search))
	for _, item := range result.Query.Search {
		results = append(results, map[string]string{
			"title": utils.NormalizeText(item.Title),
			"href":  fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(item.Title)),
			"body":  utils.NormalizeText(item.Snippet),
		})
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for: %s", query)
	}

	return results, nil
}
