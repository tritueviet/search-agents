// Package yahoo implements Yahoo search engine.
package yahoo

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

// Yahoo implements the Yahoo search engine.
type Yahoo struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	yahooSearchURL = "https://search.yahoo.com/search"
)

// New creates a new Yahoo engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &Yahoo{
		client:   client,
		name:     "yahoo",
		category: engine.CategoryText,
		provider: "yahoo",
		priority: 0.7,
	}
}

// Name returns the engine name.
func (y *Yahoo) Name() string {
	return y.name
}

// Category returns the search category.
func (y *Yahoo) Category() engine.Category {
	return y.category
}

// Provider returns the provider name.
func (y *Yahoo) Provider() string {
	return y.provider
}

// Priority returns the engine priority.
func (y *Yahoo) Priority() float64 {
	return y.priority
}

// Search performs a Yahoo search.
func (y *Yahoo) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := y.buildPayload(query, opts)

	resp, err := y.client.Do(ctx, "GET", fmt.Sprintf("%s?%s", yahooSearchURL, payload.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	results, err := y.extractResults(resp.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to extract results: %w", err)
	}

	return y.postProcessResults(results), nil
}

// buildPayload creates the search payload.
func (y *Yahoo) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("p", query)

	country, _ := parseRegion(opts.Region)
	payload.Set("gl", country)

	if opts.TimeLimit != "" {
		btfMap := map[string]string{"d": "day", "w": "week", "m": "month"}
		if btf, ok := btfMap[opts.TimeLimit]; ok {
			payload.Set("btf", btf)
		}
	}

	if opts.Page > 1 {
		payload.Set("b", fmt.Sprintf("%d", (opts.Page-1)*10+1))
	}

	return payload
}

// extractResults parses HTML and extracts search results.
func (y *Yahoo) extractResults(htmlText string) ([]models.TextResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []models.TextResult

	doc.Find("div.itm").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3.title a").Text()
		href, _ := s.Find("h3.title a").Attr("href")
		body := s.Find("div.abstract").Text()

		if title != "" {
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
func (y *Yahoo) postProcessResults(results []models.TextResult) []map[string]string {
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
