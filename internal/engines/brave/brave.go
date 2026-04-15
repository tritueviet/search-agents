// Package brave implements Brave search engine.
package brave

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

// Brave implements the Brave search engine.
type Brave struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	braveSearchURL = "https://search.brave.com/search"
)

// New creates a new Brave engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &Brave{
		client:   client,
		name:     "brave",
		category: engine.CategoryText,
		provider: "brave",
		priority: 0.8,
	}
}

// Name returns the engine name.
func (b *Brave) Name() string {
	return b.name
}

// Category returns the search category.
func (b *Brave) Category() engine.Category {
	return b.category
}

// Provider returns the provider name.
func (b *Brave) Provider() string {
	return b.provider
}

// Priority returns the engine priority.
func (b *Brave) Priority() float64 {
	return b.priority
}

// Search performs a Brave search.
func (b *Brave) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := b.buildPayload(query, opts)

	resp, err := b.client.Do(ctx, "GET", fmt.Sprintf("%s?%s", braveSearchURL, payload.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	results, err := b.extractResults(resp.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to extract results: %w", err)
	}

	return b.postProcessResults(results), nil
}

// buildPayload creates the search payload.
func (b *Brave) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("q", query)

	country, _ := parseRegion(opts.Region)
	payload.Set("country", country)

	if opts.TimeLimit != "" {
		tfMap := map[string]string{"d": "d", "w": "w", "m": "m", "y": "y"}
		if tf, ok := tfMap[opts.TimeLimit]; ok {
			payload.Set("tf", tf)
		}
	}

	if opts.Page > 1 {
		payload.Set("offset", fmt.Sprintf("%d", (opts.Page-1)*10))
	}

	return payload
}

// extractResults parses HTML and extracts search results.
func (b *Brave) extractResults(htmlText string) ([]models.TextResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []models.TextResult

	doc.Find("div.snippet").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a.title").Text()
		href, _ := s.Find("a[href]").First().Attr("href")
		body := s.Find("div.snippet-description").Text()

		if title != "" && href != "" {
			results = append(results, models.TextResult{
				Title: utils.NormalizeText(title),
				Href:  utils.NormalizeURL(href),
				Body:  utils.NormalizeText(body),
			})
		}
	})

	return results, nil
}

// postProcessResults filters empty results.
func (b *Brave) postProcessResults(results []models.TextResult) []map[string]string {
	var filtered []map[string]string

	for _, r := range results {
		if r.Title == "" || r.Href == "" {
			continue
		}
		filtered = append(filtered, map[string]string{
			"title": r.Title,
			"href":  r.Href,
			"body":  r.Body,
		})
	}

	return filtered
}

// parseRegion parses region string.
func parseRegion(region string) (string, string) {
	parts := strings.Split(strings.ToLower(region), "-")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "us", "en"
}
