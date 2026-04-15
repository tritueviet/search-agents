// Package google implements Google search engine.
package google

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

// Google implements the Google search engine.
type Google struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	googleSearchURL = "https://www.google.com/search"
)

// New creates a new Google engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &Google{
		client:   client,
		name:     "google",
		category: engine.CategoryText,
		provider: "google",
		priority: 0.85,
	}
}

// Name returns the engine name.
func (g *Google) Name() string {
	return g.name
}

// Category returns the search category.
func (g *Google) Category() engine.Category {
	return g.category
}

// Provider returns the provider name.
func (g *Google) Provider() string {
	return g.provider
}

// Priority returns the engine priority.
func (g *Google) Priority() float64 {
	return g.priority
}

// Search performs a Google search.
func (g *Google) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := g.buildPayload(query, opts)

	resp, err := g.client.Do(ctx, "GET", fmt.Sprintf("%s?%s", googleSearchURL, payload.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	results, err := g.extractResults(resp.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to extract results: %w", err)
	}

	return g.postProcessResults(results), nil
}

// buildPayload creates the search payload.
func (g *Google) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("q", query)

	country, lang := parseRegion(opts.Region)
	payload.Set("hl", lang)
	payload.Set("gl", country)

	if opts.TimeLimit != "" {
		tbsMap := map[string]string{"d": "qdr:d", "w": "qdr:w", "m": "qdr:m", "y": "qdr:y"}
		if tbs, ok := tbsMap[opts.TimeLimit]; ok {
			payload.Set("tbs", tbs)
		}
	}

	if opts.Page > 1 {
		payload.Set("start", fmt.Sprintf("%d", (opts.Page-1)*10))
	}

	return payload
}

// extractResults parses HTML and extracts search results.
func (g *Google) extractResults(htmlText string) ([]models.TextResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []models.TextResult

	doc.Find("div.g").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3").Text()
		href, _ := s.Find("a").First().Attr("href")
		body := s.Find("div.VwiC3b, div.yXK3lf").Text()

		result := models.TextResult{
			Title: utils.NormalizeText(title),
			Href:  utils.NormalizeURL(href),
			Body:  utils.NormalizeText(body),
		}

		results = append(results, result)
	})

	return results, nil
}

// postProcessResults filters out Google internal links.
func (g *Google) postProcessResults(results []models.TextResult) []map[string]string {
	var filtered []map[string]string

	for _, r := range results {
		// Skip Google internal URLs
		if strings.HasPrefix(r.Href, "/url?") || strings.HasPrefix(r.Href, "/search?") {
			continue
		}
		if r.Href == "" {
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

// parseRegion parses region string into country and language.
func parseRegion(region string) (string, string) {
	parts := strings.Split(strings.ToLower(region), "-")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "us", "en"
}
