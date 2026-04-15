// Package duckduckgo_news implements DuckDuckGo news search.
package duckduckgo_news

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/utils"
)

// DuckduckgoNews implements DuckDuckGo news search engine.
type DuckduckgoNews struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	ddgNewsURL = "https://duckduckgo.com/news.js"
	ddgMainURL = "https://duckduckgo.com"
)

// New creates a new DuckDuckGo news engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &DuckduckgoNews{
		client:   client,
		name:     "duckduckgo_news",
		category: engine.CategoryNews,
		provider: "duckduckgo",
		priority: 0.8,
	}
}

func (d *DuckduckgoNews) Name() string         { return d.name }
func (d *DuckduckgoNews) Category() engine.Category { return d.category }
func (d *DuckduckgoNews) Provider() string     { return d.provider }
func (d *DuckduckgoNews) Priority() float64    { return d.priority }

// Search performs DuckDuckGo news search.
func (d *DuckduckgoNews) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	vqd, err := d.getVQD(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get VQD: %w", err)
	}

	payload := d.buildPayload(query, opts, vqd)
	apiURL := fmt.Sprintf("%s?%s", ddgNewsURL, payload.Encode())

	resp, err := d.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	return d.extractResults(resp.Text)
}

func (d *DuckduckgoNews) getVQD(ctx context.Context, query string) (string, error) {
	resp, err := d.client.Do(ctx, "GET", ddgMainURL+"?q="+url.QueryEscape(query), nil)
	if err != nil {
		return "", err
	}

	vqd := extractVQD(resp.Text)
	if vqd == "" {
		return "", fmt.Errorf("failed to extract VQD")
	}

	return vqd, nil
}

func extractVQD(htmlText string) string {
	patterns := []string{`vqd="`, `vqd='`, `vqd=`}
	endChars := []string{`"`, `'`, `&`}

	for i, pattern := range patterns {
		idx := strings.Index(htmlText, pattern)
		if idx == -1 {
			continue
		}

		start := idx + len(pattern)
		endChar := endChars[i]
		end := strings.Index(htmlText[start:], endChar)
		if end == -1 {
			continue
		}

		return htmlText[start : start+end]
	}

	return ""
}

func (d *DuckduckgoNews) buildPayload(query string, opts engine.SearchOptions, vqd string) url.Values {
	payload := url.Values{}
	payload.Set("l", opts.Region)
	payload.Set("o", "json")
	payload.Set("q", query)
	payload.Set("vqd", vqd)

	safesearchMap := map[string]string{
		"on":       "1",
		"moderate": "-1",
		"off":      "-2",
	}
	payload.Set("p", safesearchMap[strings.ToLower(opts.SafeSearch)])

	if opts.TimeLimit != "" {
		payload.Set("df", opts.TimeLimit)
	}

	return payload
}

func (d *DuckduckgoNews) extractResults(jsonText string) ([]map[string]string, error) {
	var result struct {
		Results []struct {
			Date   string `json:"date"`
			Title  string `json:"title"`
			Excerpt string `json:"excerpt"`
			URL    string `json:"url"`
			Image  string `json:"image"`
			Source string `json:"source"`
		} `json:"results"`
	}

	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	news := make([]map[string]string, 0, len(result.Results))
	for _, r := range result.Results {
		news = append(news, map[string]string{
			"date":   r.Date,
			"title":  utils.NormalizeText(r.Title),
			"body":   utils.NormalizeText(r.Excerpt),
			"url":    utils.NormalizeURL(r.URL),
			"image":  utils.NormalizeURL(r.Image),
			"source": utils.NormalizeText(r.Source),
		})
	}

	if len(news) == 0 {
		return nil, fmt.Errorf("no news extracted")
	}

	return news[:min(len(news), 10)], nil
}
