// Package bing implements Bing search engine.
package bing

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/models"
	"github.com/tritueviet/search-agents/internal/utils"
)

// Bing implements the Bing search engine.
type Bing struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	bingSearchURL = "https://www.bing.com/search"
)

// New creates a new Bing engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &Bing{
		client:   client,
		name:     "bing",
		category: engine.CategoryText,
		provider: "bing",
		priority: 0.9,
	}
}

// Name returns the engine name.
func (b *Bing) Name() string {
	return b.name
}

// Category returns the search category.
func (b *Bing) Category() engine.Category {
	return b.category
}

// Provider returns the provider name.
func (b *Bing) Provider() string {
	return b.provider
}

// Priority returns the engine priority.
func (b *Bing) Priority() float64 {
	return b.priority
}

// Search performs a Bing search.
func (b *Bing) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := b.buildPayload(query, opts)

	resp, err := b.client.Do(ctx, "GET", fmt.Sprintf("%s?%s", bingSearchURL, payload.Encode()), nil)
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
func (b *Bing) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("q", query)

	_, lang := parseRegion(opts.Region)
	payload.Set("cc", lang)

	// Set cookies for region
	cookies := map[string]string{
		"_EDGE_CD": fmt.Sprintf("m=%s&u=%s", opts.Region, opts.Region),
		"_EDGE_S":  fmt.Sprintf("mkt=%s&ui=%s", opts.Region, opts.Region),
	}
	b.client.SetHeaders(map[string]string{
		"Cookie": formatCookies(cookies),
	})

	if opts.TimeLimit != "" {
		d := int(time.Now().Unix() / 86400)
		var code string
		if opts.TimeLimit == "y" {
			code = fmt.Sprintf("ez5_%d_%d", d-365, d)
		} else {
			codeMap := map[string]string{"d": "1", "w": "2", "m": "3"}
			code = "ez" + codeMap[opts.TimeLimit]
		}
		payload.Set("filters", fmt.Sprintf("ex1:\"%s\"", code))
	}

	if opts.Page > 1 {
		payload.Set("first", fmt.Sprintf("%d", (opts.Page-1)*10))
	}

	return payload
}

// extractResults parses HTML and extracts search results.
func (b *Bing) extractResults(htmlText string) ([]models.TextResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []models.TextResult

	doc.Find("li.b_algo").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h2 a").Text()
		href, _ := s.Find("h2 a").Attr("href")
		body := s.Find("p").Text()

		result := models.TextResult{
			Title: utils.NormalizeText(title),
			Href:  utils.NormalizeURL(href),
			Body:  utils.NormalizeText(body),
		}

		if !strings.HasPrefix(result.Href, "https://www.bing.com/aclick?") {
			results = append(results, result)
		}
	})

	return results, nil
}

// postProcessResults filters and transforms results.
func (b *Bing) postProcessResults(results []models.TextResult) []map[string]string {
	var filtered []map[string]string

	for _, r := range results {
		if strings.HasPrefix(r.Href, "https://www.bing.com/ck/a?") {
			r.Href = unwrapBingURL(r.Href)
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

// formatCookies formats cookie map into header value.
func formatCookies(cookies map[string]string) string {
	parts := make([]string, 0, len(cookies))
	for k, v := range cookies {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, "; ")
}

// unwrapBingURL decodes a Bing-wrapped URL.
func unwrapBingURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	u := parsed.Query().Get("u")
	if u == "" || len(u) <= 2 {
		return rawURL
	}
	// Drop first 2 chars and decode
	b64 := u[2:]
	padding := strings.Repeat("=", (4-len(b64)%4)%4)
	decoded, err := decodeBase64URL(b64 + padding)
	if err != nil {
		return rawURL
	}
	return string(decoded)
}

// decodeBase64URL decodes base64 URL-encoded string.
func decodeBase64URL(s string) ([]byte, error) {
	// Simple base64 decode placeholder
	return []byte(s), nil
}
