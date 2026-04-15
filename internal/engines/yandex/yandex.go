// Package yandex implements Yandex search engine.
package yandex

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

// Yandex implements the Yandex search engine.
type Yandex struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	yandexSearchURL = "https://yandex.com/search"
)

// New creates a new Yandex engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &Yandex{
		client:   client,
		name:     "yandex",
		category: engine.CategoryText,
		provider: "yandex",
		priority: 0.75,
	}
}

// Name returns the engine name.
func (y *Yandex) Name() string {
	return y.name
}

// Category returns the search category.
func (y *Yandex) Category() engine.Category {
	return y.category
}

// Provider returns the provider name.
func (y *Yandex) Provider() string {
	return y.provider
}

// Priority returns the engine priority.
func (y *Yandex) Priority() float64 {
	return y.priority
}

// Search performs a Yandex search.
func (y *Yandex) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := y.buildPayload(query, opts)

	resp, err := y.client.Do(ctx, "GET", fmt.Sprintf("%s?%s", yandexSearchURL, payload.Encode()), nil)
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
func (y *Yandex) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("text", query)

	_, lang := parseRegion(opts.Region)
	payload.Set("lr", fmt.Sprintf("lang:%s", lang))

	if opts.TimeLimit != "" {
		withinMap := map[string]string{"d": "1", "w": "2", "m": "3", "y": "4"}
		if within, ok := withinMap[opts.TimeLimit]; ok {
			payload.Set("within", within)
		}
	}

	if opts.Page > 1 {
		payload.Set("p", fmt.Sprintf("%d", opts.Page-1))
	}

	return payload
}

// extractResults parses HTML and extracts search results.
func (y *Yandex) extractResults(htmlText string) ([]models.TextResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []models.TextResult

	doc.Find("li.serp-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h2 a").Text()
		href, _ := s.Find("h2 a").Attr("href")
		body := s.Find("div.text-container").Text()

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
func (y *Yandex) postProcessResults(results []models.TextResult) []map[string]string {
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
