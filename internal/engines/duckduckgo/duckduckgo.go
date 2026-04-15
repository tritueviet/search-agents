// Package duckduckgo implements the DuckDuckGo search engine.
package duckduckgo

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/models"
	"github.com/tritueviet/search-agents/internal/utils"
)

// DuckduckGo implements the DuckDuckGo search engine.
type DuckduckGo struct {
	client     *httpclient.Client
	name       string
	category   engine.Category
	provider   string
	priority   float64
}

const (
	ddgSearchURL = "https://html.duckduckgo.com/html/"
)

// New creates a new DuckDuckGo engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &DuckduckGo{
		client:   client,
		name:     "duckduckgo",
		category: engine.CategoryText,
		provider: "bing",
		priority: 1.0,
	}
}

// Name returns the engine name.
func (d *DuckduckGo) Name() string {
	return d.name
}

// Category returns the search category.
func (d *DuckduckGo) Category() engine.Category {
	return d.category
}

// Provider returns the provider name.
func (d *DuckduckGo) Provider() string {
	return d.provider
}

// Priority returns the engine priority.
func (d *DuckduckGo) Priority() float64 {
	return d.priority
}

// Search performs a DuckDuckGo search.
func (d *DuckduckGo) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := d.buildPayload(query, opts)

	resp, err := d.client.PostForm(ctx, ddgSearchURL, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d, body: %s", resp.StatusCode, resp.Text[:200])
	}

	results, err := d.extractResults(resp.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to extract results: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results extracted from HTML (status %d, body length %d)", resp.StatusCode, len(resp.Text))
	}

	return d.postProcessResults(results), nil
}

// buildPayload creates the search payload.
func (d *DuckduckGo) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("q", query)
	payload.Set("b", "")
	payload.Set("l", opts.Region)

	if opts.Page > 1 {
		payload.Set("s", fmt.Sprintf("%d", 10+(opts.Page-2)*15))
	}
	if opts.TimeLimit != "" {
		payload.Set("df", opts.TimeLimit)
	}

	return payload
}

// extractResults parses HTML and extracts search results.
func (d *DuckduckGo) extractResults(htmlText string) ([]models.TextResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []models.TextResult

	// Try multiple selectors for DuckDuckGo
	// Primary: div.web-result (new structure)
	// Fallback: div.result, div.body (old structure)
	selectors := []string{
		"div.web-result",
		"div.result",
		"div.body",
		"article[data-testid='result']",
	}

	var items *goquery.Selection
	for _, selector := range selectors {
		items = doc.Find(selector)
		if items.Length() > 0 {
			break
		}
	}

	if items.Length() == 0 {
		return nil, fmt.Errorf("no results found with any selector")
	}

	items.Each(func(i int, s *goquery.Selection) {
		// Try multiple title selectors
		title := s.Find("h2 a").Text()
		if title == "" {
			title = s.Find("a.result__a").Text()
		}
		if title == "" {
			title = s.Find("h2").Text()
		}

		// Try multiple href selectors
		href, _ := s.Find("h2 a").Attr("href")
		if href == "" {
			href, _ = s.Find("a.result__a").Attr("href")
		}
		if href == "" {
			href, _ = s.Find("a").First().Attr("href")
		}

		// Try multiple body selectors
		body := s.Find("a").Text()
		if body == "" {
			body = s.Find("p").Text()
		}
		if body == "" {
			body = s.Find("div.result__snippet").Text()
		}
		if body == "" {
			body = s.Find("div.result__body").Text()
		}

		result := models.TextResult{
			Title: utils.NormalizeText(title),
			Href:  utils.NormalizeURL(href),
			Body:  utils.NormalizeText(body),
		}

		if result.Title != "" || result.Href != "" {
			results = append(results, result)
		}
	})

	return results, nil
}

// postProcessResults filters out unwanted results.
func (d *DuckduckGo) postProcessResults(results []models.TextResult) []map[string]string {
	var filtered []map[string]string

	for _, r := range results {
		if !strings.HasPrefix(r.Href, "https://duckduckgo.com/y.js?") {
			filtered = append(filtered, map[string]string{
				"title": r.Title,
				"href":  r.Href,
				"body":  r.Body,
			})
		}
	}

	return filtered
}
